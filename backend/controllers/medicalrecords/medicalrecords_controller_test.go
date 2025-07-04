package medicalrecords_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/medicalrecords"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	TestAdminUser   models.User
	TestDoctorUser  models.User
	TestPatientUser models.User
)

func init() {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.ConnectDB()

	TestAdminUser = seedUser(models.RoleAdmin)
	TestDoctorUser = seedUser(models.RoleDoctor)
	TestPatientUser = seedUser(models.RolePatient)

	config.DB.FirstOrCreate(&models.Doctor{ID: TestDoctorUser.ID, UserID: TestDoctorUser.ID})
	config.DB.FirstOrCreate(&models.Patient{ID: TestPatientUser.ID, UserID: TestPatientUser.ID})
}

func seedUser(role models.Role) models.User {
	email := fmt.Sprintf("test_%s@example.com", role)
	user := models.User{}
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}
	user = models.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  string(role),
		Email:     email,
		Password:  "securepass",
		Role:      role,
	}
	config.DB.Create(&user)
	return user
}

func injectTestUser(user models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwtPayload", &user)
		c.Set("user", &user)
		c.Next()
	}
}

func testRoleChecker(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("jwtPayload")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}
		log.Printf("val type: %T", val)
		userPtr, ok := val.(*models.User)
		if !ok || userPtr == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
			c.Abort()
			return
		}

		user := *userPtr

		if slices.Contains(allowedRoles, user.Role) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Forbidden: role '%s' not permitted", user.Role)})
		c.Abort()
	}
}

func setupMedicalRecordTestRouter(user models.User) *gin.Engine {
	r := gin.Default()
	r.Use(injectTestUser(user))

	rg := r.Group("/medicalrecords")
	rg.POST("/", testRoleChecker(models.RoleDoctor, models.RoleAdmin), medicalrecords.CreateMedicalRecord)
	rg.GET("/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), medicalrecords.GetMedicalRecord)
	rg.GET("/patient/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), medicalrecords.GetRecordsByPatientID)
	rg.PUT("/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), medicalrecords.UpdateMedicalRecord)
	rg.DELETE("/soft/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), medicalrecords.SoftDeleteMedicalRecord)
	rg.DELETE("/hard/:id", testRoleChecker(models.RoleAdmin), medicalrecords.HardDeleteMedicalRecord)

	return r
}

func createTestMedicalRecord(t *testing.T) models.MedicalRecord {
	record := models.MedicalRecord{
		ID:        uuid.New(),
		PatientID: TestPatientUser.ID,
		DoctorID:  TestDoctorUser.ID,
		Diagnosis: "Test Diagnosis",
		Notes:     "Initial notes",
	}
	assert.NoError(t, config.DB.Create(&record).Error)
	return record
}

func TestCreateMedicalRecord(t *testing.T) {
	router := setupMedicalRecordTestRouter(TestDoctorUser)

	// --- Case 1: With Vitals ---
	vitals := []map[string]interface{}{
		{
			"type":        "BloodPressure",
			"value":       "120/80",
			"status":      "Normal",
			"recorded_at": "2025-07-04T10:00:00Z",
		},
		{
			"type":        "HeartRate",
			"value":       "75",
			"status":      "Normal",
			"recorded_at": "2025-07-04T10:10:00Z",
		},
	}
	withVitalsBody := map[string]interface{}{
		"patient_id": TestPatientUser.ID,
		"doctor_id":  TestDoctorUser.ID,
		"diagnosis":  "High BP observed",
		"notes":      "Prescribed mild medication",
		"vitals":     vitals,
	}
	withVitalsJSON, _ := json.Marshal(withVitalsBody)

	req1 := httptest.NewRequest("POST", "/medicalrecords/", bytes.NewBuffer(withVitalsJSON))
	req1.Header.Set("Content-Type", "application/json")
	res1 := httptest.NewRecorder()
	router.ServeHTTP(res1, req1)

	assert.Equal(t, http.StatusCreated, res1.Code)
	assert.Contains(t, res1.Body.String(), "Medical record created successfully")

	// --- Case 2: Without Vitals ---
	withoutVitalsBody := map[string]interface{}{
		"patient_id": TestPatientUser.ID,
		"doctor_id":  TestDoctorUser.ID,
		"diagnosis":  "Follow-up visit",
		"notes":      "No vitals recorded today",
	}
	withoutVitalsJSON, _ := json.Marshal(withoutVitalsBody)

	req2 := httptest.NewRequest("POST", "/medicalrecords/", bytes.NewBuffer(withoutVitalsJSON))
	req2.Header.Set("Content-Type", "application/json")
	res2 := httptest.NewRecorder()
	router.ServeHTTP(res2, req2)

	assert.Equal(t, http.StatusCreated, res2.Code)
	assert.Contains(t, res2.Body.String(), "Medical record created successfully")
}

func TestGetMedicalRecord(t *testing.T) {
	r := setupMedicalRecordTestRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	req := httptest.NewRequest("GET", "/medicalrecords/"+record.ID.String(), nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Test Diagnosis")
}

func TestGetRecordsByPatientID(t *testing.T) {
	r := setupMedicalRecordTestRouter(TestDoctorUser)
	_ = createTestMedicalRecord(t)

	req := httptest.NewRequest("GET", "/medicalrecords/patient/"+TestPatientUser.ID.String(), nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Test Diagnosis")
}

func TestUpdateMedicalRecord(t *testing.T) {
	r := setupMedicalRecordTestRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	update := map[string]string{
		"diagnosis": "Updated Diagnosis",
		"notes":     "Updated notes",
	}
	jsonUpdate, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/medicalrecords/"+record.ID.String(), bytes.NewBuffer(jsonUpdate))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Medical record updated")
}

func TestSoftDeleteMedicalRecord(t *testing.T) {
	r := setupMedicalRecordTestRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	req := httptest.NewRequest("DELETE", "/medicalrecords/soft/"+record.ID.String(), nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "soft deleted")
}

func TestHardDeleteMedicalRecord(t *testing.T) {
	r := setupMedicalRecordTestRouter(TestAdminUser)
	record := createTestMedicalRecord(t)

	req := httptest.NewRequest("DELETE", "/medicalrecords/hard/"+record.ID.String(), nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "hard deleted")
}

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
	password := "securepass"

	// Check if Auth already exists
	var auth models.Auth
	if err := config.DB.Where("email = ?", email).First(&auth).Error; err != nil {
		auth = models.Auth{
			ID:       uuid.New(),
			Email:    email,
			Password: password, // You might want to hash this in real app logic
		}
		config.DB.Create(&auth)
	}

	// Check if User already exists for that Auth
	var user models.User
	if err := config.DB.Where("auth_id = ?", auth.ID).First(&user).Error; err != nil {
		user = models.User{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  string(role),
			AuthID:    auth.ID,
			Role:      role,
			Phone:     "1234567890",
		}
		config.DB.Create(&user)
	}
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
		user, ok := val.(*models.User)
		if !ok || user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
			c.Abort()
			return
		}
		if slices.Contains(allowedRoles, user.Role) {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Forbidden: role '%s' not permitted", user.Role)})
		c.Abort()
	}
}

func setupRouter(user models.User) *gin.Engine {
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

// ---------------------- TEST CASES ----------------------

func TestCreateMedicalRecord(t *testing.T) {
	router := setupRouter(TestDoctorUser)

	t.Run("Create With Vitals", func(t *testing.T) {
		body := map[string]interface{}{
			"patient_id": TestPatientUser.ID,
			"doctor_id":  TestDoctorUser.ID,
			"diagnosis":  "High BP",
			"notes":      "Prescribed medication",
			"vitals": []map[string]interface{}{
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
			},
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/medicalrecords/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), "Medical record created successfully")
	})

	t.Run("Create Without Vitals", func(t *testing.T) {
		body := map[string]interface{}{
			"patient_id": TestPatientUser.ID,
			"doctor_id":  TestDoctorUser.ID,
			"diagnosis":  "Routine checkup",
			"notes":      "All good",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/medicalrecords/", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), "Medical record created successfully")
	})
}

func TestGetMedicalRecord(t *testing.T) {
	r := setupRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	t.Run("Fetch Medical Record By ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/medicalrecords/"+record.ID.String(), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Test Diagnosis")
	})
}

func TestGetRecordsByPatientID(t *testing.T) {
	r := setupRouter(TestDoctorUser)
	_ = createTestMedicalRecord(t)

	t.Run("Fetch Records By Patient ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/medicalrecords/patient/"+TestPatientUser.ID.String(), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Test Diagnosis")
	})
}

func TestUpdateMedicalRecord(t *testing.T) {
	r := setupRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	t.Run("Update Diagnosis and Notes", func(t *testing.T) {
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
	})
}

func TestSoftDeleteMedicalRecord(t *testing.T) {
	r := setupRouter(TestDoctorUser)
	record := createTestMedicalRecord(t)

	t.Run("Soft Delete Record", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/medicalrecords/soft/"+record.ID.String(), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "soft deleted")
	})
}

func TestHardDeleteMedicalRecord(t *testing.T) {
	r := setupRouter(TestAdminUser)
	record := createTestMedicalRecord(t)

	t.Run("Hard Delete Record", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/medicalrecords/hard/"+record.ID.String(), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "hard deleted")
	})
}

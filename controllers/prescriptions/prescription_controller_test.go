package prescriptions_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/prescriptions"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	TestPatient models.User
	TestDoctor  models.User
	TestAdmin   models.User
)

func init() {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()
	_ = godotenv.Load("../../.env")
	config.ConnectDB()

	TestPatient = seedPrescriptionUser(models.RolePatient)
	TestDoctor = seedPrescriptionUser(models.RoleDoctor)
	TestAdmin = seedPrescriptionUser(models.RoleAdmin)
}

func seedPrescriptionUser(role models.Role) models.User {
	email := string(role) + "_rx_test@example.com"
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}
	user = models.User{
		ID:        uuid.New(),
		FirstName: "Rx",
		LastName:  "Tester",
		Role:      role,
	}
	config.DB.Create(&user)
	switch role {
	case models.RolePatient:
		config.DB.FirstOrCreate(&models.Patient{ID: user.ID, UserID: user.ID})
	case models.RoleDoctor:
		config.DB.FirstOrCreate(&models.Doctor{ID: user.ID, UserID: user.ID})
	}
	return user
}

func injectPrescriptionUser(user models.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwtPayload", &user)
		c.Set("user", &user)
		c.Next()
	}
}

func setupRxRouter(user models.User) *gin.Engine {
	r := gin.Default()
	r.Use(injectPrescriptionUser(user))

	r.POST("/prescriptions", prescriptions.CreatePrescription)
	r.GET("/prescriptions/patient/:id", prescriptions.GetPrescriptionsByPatientID)
	r.GET("/prescriptions/:id", prescriptions.GetPrescriptionByID)
	r.PUT("/prescriptions/:id", prescriptions.UpdatePrescription)
	r.DELETE("/prescriptions/:id", prescriptions.DeletePrescription)
	return r
}

func createTestMedicalRecord(t *testing.T, patientID, doctorID uuid.UUID) uuid.UUID {
	record := models.MedicalRecord{
		ID:        uuid.New(),
		PatientID: patientID,
		DoctorID:  doctorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := config.DB.Create(&record).Error; err != nil {
		t.Fatalf("Failed to create test medical record: %v", err)
	}
	return record.ID
}

func createTestPrescription(t *testing.T, router *gin.Engine, patientID, doctorID uuid.UUID) uuid.UUID {
	medicalRecordID := createTestMedicalRecord(t, patientID, doctorID)

	input := map[string]interface{}{
		"patient_id":        patientID,
		"doctor_id":         doctorID,
		"medical_record_id": medicalRecordID,
		"medication":        "TestMed",
		"dosage":            "1x daily",
		"instructions":      "After food",
		"issued_at":         time.Now().Format(time.RFC3339),
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/prescriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	idStr, ok := res["prescription_id"].(string)
	if !ok {
		t.Fatalf("Expected string prescription_id, got: %v", res["prescription_id"])
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		t.Fatalf("Invalid UUID: %v", err)
	}
	return id
}

func TestCreatePrescription(t *testing.T) {
	router := setupRxRouter(TestDoctor)
	_ = createTestPrescription(t, router, TestPatient.ID, TestDoctor.ID)
}

func TestGetPrescriptionByID(t *testing.T) {
	router := setupRxRouter(TestDoctor)
	id := createTestPrescription(t, router, TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("GET", "/prescriptions/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestMed")
}

func TestGetPrescriptionsByPatientID(t *testing.T) {
	router := setupRxRouter(TestPatient)
	_ = createTestPrescription(t, setupRxRouter(TestDoctor), TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("GET", "/prescriptions/patient/"+TestPatient.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestMed")
}

func TestUpdatePrescription(t *testing.T) {
	router := setupRxRouter(TestDoctor)
	id := createTestPrescription(t, router, TestPatient.ID, TestDoctor.ID)

	update := map[string]interface{}{
		"dosage":    "2x daily",
		"issued_at": time.Now().Format(time.RFC3339),
	}
	b, _ := json.Marshal(update)
	req := httptest.NewRequest("PUT", "/prescriptions/"+id.String(), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestDeletePrescription(t *testing.T) {
	router := setupRxRouter(TestDoctor)
	id := createTestPrescription(t, router, TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("DELETE", "/prescriptions/"+id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted")
}

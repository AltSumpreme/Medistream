package appointments_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/appointments"
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
	email := fmt.Sprintf("test_%s_%s@example.com", role, uuid.New().String())

	// Check if user already exists
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}

	// Create Auth entry first
	auth := models.Auth{
		ID:       uuid.New(),
		Email:    email,
		Password: "hashed_dummy_password", // In production, hash it properly
	}
	if err := config.DB.Create(&auth).Error; err != nil {
		log.Fatalf("failed to seed auth: %v", err)
	}

	user = models.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  string(role),
		Phone:     "1234567890",
		Role:      role,
		AuthID:    auth.ID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Fatalf("failed to seed user: %v", err)
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

func setupRouterWithTestRoutes(user models.User) *gin.Engine {
	r := gin.Default()
	r.Use(injectTestUser(user))
	r.Use(gin.Logger())

	rg := r.Group("/appointments")
	rg.POST("/", testRoleChecker(models.RolePatient, models.RoleDoctor, models.RoleAdmin), appointments.CreateAppointment)
	rg.GET("/", testRoleChecker(models.RoleAdmin), appointments.GetAllAppointments)
	rg.GET("/:id", testRoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), appointments.GetAppointmentByID)
	rg.PUT("/:id", testRoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), appointments.UpdateAppointment)
	rg.PUT("/status/:id", testRoleChecker(models.RoleAdmin), appointments.ChangeAppointmentStatus)
	rg.PUT("/reschedule/:id", testRoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.RescheduleAppointment)
	rg.PUT("/cancel/:id", testRoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.CancelAppointment)
	rg.GET("/doctor/:id", testRoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.GetAppointmentByDoctorID)
	rg.GET("/patient/:id", testRoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.GetAppointmentByPatientID)
	rg.DELETE("/:id", testRoleChecker(models.RoleAdmin), appointments.DeleteAppointment)

	return r
}

func createTestAppointment(patientID, doctorID uuid.UUID) models.Appointment {
	appt := models.Appointment{
		ID:              uuid.New(),
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentDate: time.Now().Add(100 * time.Hour),
		Duration:        30,
		Location:        "Test Room",
		Mode:            "Online",
		AppointmentType: models.ApptType("Consultation"),
		Status:          models.AppointmentStatusPending,
	}
	config.DB.Create(&appt)
	return appt
}

// -------------------- TESTS -----------------------

func TestCreateAppointment(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	body := map[string]interface{}{
		"patient_id":       TestPatientUser.ID,
		"doctor_id":        TestDoctorUser.ID,
		"appointment_date": time.Now().Add(200 * time.Hour).Format(time.RFC3339),
		"duration":         45,
		"appointment_type": "CONSULTATION",
		"location":         "Room A",
		"mode":             "Online",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/appointments/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment created successfully")
}

func TestGetAllAppointments(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)

	req := httptest.NewRequest("GET", "/appointments/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByID(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	req := httptest.NewRequest("GET", "/appointments/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointment")
}

func TestUpdateAppointment(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	update := map[string]interface{}{
		"patient_id":       TestPatientUser.ID,
		"doctor_id":        TestDoctorUser.ID,
		"appointment_date": time.Now().Add(72 * time.Hour).Format(time.RFC3339),
		"duration":         60,
		"notes":            "Updated notes",
		"appointment_type": "FOLLOWUP",
		"location":         "Updated Room",
		"mode":             "In-Person",
	}
	jsonBody, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/appointments/"+appt.ID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment updated successfully")
}

func TestDeleteAppointment(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	req := httptest.NewRequest("DELETE", "/appointments/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment deleted successfully")
}

func TestCancelAppointment(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	req := httptest.NewRequest("PUT", "/appointments/cancel/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment cancelled")
}

func TestRescheduleAppointment(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	body := map[string]interface{}{
		"date":     time.Now().Add(96 * time.Hour).Format(time.RFC3339),
		"duration": 60,
		"mode":     "In-Person",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/appointments/reschedule/"+appt.ID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment rescheduled")
}

func TestChangeAppointmentStatus(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	appt := createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	body := map[string]interface{}{
		"status": "COMPLETED",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/appointments/status/"+appt.ID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetAppointmentByDoctorID(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	_ = createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	req := httptest.NewRequest("GET", "/appointments/doctor/"+TestDoctorUser.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByPatientID(t *testing.T) {
	router := setupRouterWithTestRoutes(TestAdminUser)
	_ = createTestAppointment(TestPatientUser.ID, TestDoctorUser.ID)

	req := httptest.NewRequest("GET", "/appointments/patient/"+TestPatientUser.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

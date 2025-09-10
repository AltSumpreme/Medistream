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
	TestDoctor      models.Doctor
	TestPatient     models.Patient
)

func init() {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.ConnectDB()

	// Seed base users
	TestAdminUser = seedUser(models.RoleAdmin)
	TestDoctorUser = seedUser(models.RoleDoctor)
	TestPatientUser = seedUser(models.RolePatient)

	// Seed doctor with independent ID
	TestDoctor = models.Doctor{
		ID:             uuid.New(),
		UserID:         TestDoctorUser.ID,
		Specialization: "General",
	}
	if err := config.DB.FirstOrCreate(&TestDoctor, models.Doctor{UserID: TestDoctorUser.ID}).Error; err != nil {
		log.Fatalf("failed to seed doctor: %v", err)
	}

	// Seed patient with independent ID
	TestPatient = models.Patient{
		ID:     uuid.New(),
		UserID: TestPatientUser.ID,
	}
	if err := config.DB.FirstOrCreate(&TestPatient, models.Patient{UserID: TestPatientUser.ID}).Error; err != nil {
		log.Fatalf("failed to seed patient: %v", err)
	}
}

func seedUser(role models.Role) models.User {
	email := fmt.Sprintf("test_%s_%s@example.com", role, uuid.New().String())

	// Check if user already exists
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}

	auth := models.Auth{
		ID:       uuid.New(),
		Email:    email,
		Password: "hashed_dummy_password",
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

func injectTestUser(claims *utils.JWTClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("jwtPayload", claims)
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

		claims, ok := val.(*utils.JWTClaims)
		if !ok || claims == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
			c.Abort()
			return
		}

		if slices.Contains(allowedRoles, models.Role(claims.Role)) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Forbidden: role '%s' not permitted", claims.Role)})
		c.Abort()
	}
}

func setupRouterWithTestRoutes(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(injectTestUser(claims))

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
		AppointmentDate: time.Now().Add(48 * time.Hour),
		StartTime:       "10:00",
		EndTime:         "10:30",
		Notes:           "Initial Test Appointment",
		Location:        "Test Room",
		Mode:            "Online",
		Status:          models.AppointmentStatus("PENDING"),
		AppointmentType: models.ApptType("CONSULTATION"),
	}
	if err := config.DB.Create(&appt).Error; err != nil {
		log.Fatalf("failed to create test appointment: %v", err)
	}
	return appt
}

// -------------------- TESTS -----------------------

func TestCreateAppointment(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestPatientUser.ID,
		Role:   string(models.RolePatient),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	body := map[string]interface{}{
		"doctorId":        TestDoctor.ID,
		"appointmentDate": time.Now().Add(200 * time.Hour).Format(time.RFC3339),
		"startTime":       "14:00",
		"endTime":         "14:30",
		"appointmentType": "CONSULTATION",
		"mode":            "Online",
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
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)

	req := httptest.NewRequest("GET", "/appointments/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByID(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("GET", "/appointments/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointment")
}

func TestUpdateAppointment(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

	update := map[string]interface{}{
		"patient_id":       TestPatientUser.ID,
		"doctor_id":        TestDoctorUser.ID,
		"appointment_date": time.Now().Add(72 * time.Hour).Format(time.RFC3339),
		"start_time":       "11:00",
		"end_time":         "11:30",
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
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("DELETE", "/appointments/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment deleted successfully")
}

func TestCancelAppointment(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("PUT", "/appointments/cancel/"+appt.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment cancelled")
}

func TestRescheduleAppointment(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

	body := map[string]interface{}{
		"date":       time.Now().Add(96 * time.Hour).Format(time.RFC3339),
		"start_time": "4:00",
		"end_time":   "4:30",
		"mode":       "In-Person",
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
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	appt := createTestAppointment(TestPatient.ID, TestDoctor.ID)

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
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	_ = createTestAppointment(TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("GET", "/appointments/doctor/"+TestDoctor.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByPatientID(t *testing.T) {
	claims := &utils.JWTClaims{
		UserID: TestAdminUser.ID,
		Role:   string(models.RoleAdmin),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	router := setupRouterWithTestRoutes(claims)
	_ = createTestAppointment(TestPatient.ID, TestDoctor.ID)

	req := httptest.NewRequest("GET", "/appointments/patient/"+TestPatient.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

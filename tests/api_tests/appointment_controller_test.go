package apitests

import (
	"net/http"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	apiclient "github.com/AltSumpreme/Medistream.git/tests/api_client"
	"github.com/AltSumpreme/Medistream.git/tests/factories"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupApptRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	//r.Use(middleware.RequestTimer())
	// appointmentCache := cache.NewCache(config.Rdb, config.Ctx)
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	// rg := r.Group("/appointments")
	// routes.RegisterAppointmentRoutes(rg, appointmentCache)
	return r
}

func TestCreateAppointment(t *testing.T) {
	db := config.DB
	userPatient, _, _, doctor, _ := factories.CreateEntries(db)
	claims := factories.MakeJWT(userPatient.ID, models.RolePatient)

	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	body := map[string]interface{}{
		"doctorId":        doctor.ID,
		"appointmentDate": time.Now().Add(200 * time.Hour).Format(time.RFC3339),
		"startTime":       "14:00",
		"endTime":         "14:30",
		"appointmentType": "CONSULTATION",
		"mode":            "Online",
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Post("/appointments", body, headers)
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment created successfully")
}

func TestAppointmentCacheOnCreate(t *testing.T) {
	db := config.DB
	userPatient, _, _, doctor, _ := factories.CreateEntries(db)
	claims := factories.MakeJWT(userPatient.ID, models.RolePatient)

	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	body := map[string]interface{}{
		"doctorId":        doctor.ID,
		"appointmentDate": time.Now().Add(200 * time.Hour).Format(time.RFC3339),
		"startTime":       "14:00",
		"endTime":         "14:30",
		"appointmentType": "CONSULTATION",
		"mode":            "Online",
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Post("/appointments", body, headers)
	assert.Equal(t, http.StatusCreated, res.Code)

	// 🔑 Check Redis cache
	key := "appointments:doctor:" + doctor.ID.String()
	cached, err := config.Rdb.Get(config.Ctx, key).Result()
	assert.NoError(t, err)
	assert.Contains(t, cached, doctor.ID.String())
}

func TestGetAllAppointments(t *testing.T) {
	db := config.DB
	admin := factories.SeedUser(db, models.RoleAdmin)
	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)

	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Get("/appointments", nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByID(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	appt := factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Get("/appointments/"+appt.ID.String(), nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointment")
}

func TestUpdateAppointment(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	appt := factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	update := map[string]interface{}{
		"notes":            "Updated test notes",
		"appointmentType":  "FOLLOWUP",
		"startTime":        "15:00",
		"endTime":          "15:30",
		"appointment_date": time.Now().Add(72 * time.Hour).Format(time.RFC3339),
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Put("/appointments/"+appt.ID.String(), update, headers)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment updated successfully")
}

func TestDeleteAppointment(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	appt := factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Delete("/appointments/"+appt.ID.String(), nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment deleted successfully")
}

func TestCancelAppointment(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	appt := factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Put("/appointments/cancel/"+appt.ID.String(), nil, nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment cancelled")
}

func TestRescheduleAppointment(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	appt := factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	body := map[string]interface{}{
		"date":       time.Now().Add(96 * time.Hour).Format(time.RFC3339),
		"start_time": "16:00",
		"end_time":   "16:30",
		"mode":       "In-Person",
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Put("/appointments/reschedule/"+appt.ID.String(), body, headers)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment rescheduled")
}

func TestGetAppointmentByDoctorID(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	_ = factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Get("/appointments/doctor/"+doctor.ID.String(), nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

func TestGetAppointmentByPatientID(t *testing.T) {
	db := config.DB
	_, patient, _, doctor, _ := factories.CreateEntries(db)
	admin := factories.SeedUser(db, models.RoleAdmin)
	_ = factories.CreateAppointment(db, patient.ID, doctor.ID)

	claims := factories.MakeJWT(admin.ID, models.RoleAdmin)
	router := setupApptRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	res := client.Get("/appointments/patient/"+patient.ID.String(), nil)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "appointments")
}

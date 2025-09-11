package apitests

import (
	"net/http"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/appointments"
	"github.com/AltSumpreme/Medistream.git/models"
	apiclient "github.com/AltSumpreme/Medistream.git/tests/api_client"
	"github.com/AltSumpreme/Medistream.git/tests/factories"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	rg := r.Group("/appointments")
	rg.POST("/", appointments.CreateAppointment)
	rg.GET("/", appointments.GetAllAppointments)
	rg.GET("/:id", appointments.GetAppointmentByID)
	rg.PUT("/:id", appointments.UpdateAppointment)
	rg.PUT("/status/:id", appointments.ChangeAppointmentStatus)
	rg.PUT("/reschedule/:id", appointments.RescheduleAppointment)
	rg.PUT("/cancel/:id", appointments.CancelAppointment)
	rg.GET("/doctor/:id", appointments.GetAppointmentByDoctorID)
	rg.GET("/patient/:id", appointments.GetAppointmentByPatientID)
	rg.DELETE("/:id", appointments.DeleteAppointment)
	return r
}

func TestCreateAppointment(t *testing.T) {
	if config.DB == nil {
		t.Fatal("Database connection is not initialized")
	}
	db := config.DB
	Userpatient := factories.SeedUser(db, models.RolePatient)
	factories.SeedPatient(db, Userpatient)
	doctor := factories.SeedUser(db, models.RoleDoctor)
	doc := factories.SeedDoctor(db, doctor)
	token := factories.GenerateJWT(Userpatient.ID.String(), string(models.RolePatient))
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	router := setupRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	body := map[string]interface{}{
		"doctorId":        doc.ID,
		"appointmentDate": time.Now().Add(200 * time.Hour).Format(time.RFC3339),
		"startTime":       "14:00",
		"endTime":         "14:30",
		"appointmentType": "CONSULTATION",
		"mode":            "Online",
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res := client.Post("/appointments/", body, headers)

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, res.Body.String(), "Appointment created successfully")
}

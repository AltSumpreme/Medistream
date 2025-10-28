package apitests

import (
	"net/http"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/routes"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	apiclient "github.com/AltSumpreme/Medistream.git/tests/api_client"
	"github.com/AltSumpreme/Medistream.git/tests/factories"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupPrescriptionRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	prescriptionCache := cache.NewCache(config.Rdb, config.Ctx)
	rg := r.Group("/prescriptions")
	routes.RegisterPrescriptionRoutes(rg, prescriptionCache)
	return r
}

func TestPrescriptionRoutes(t *testing.T) {
	db := config.DB

	// Seed users
	userPatient, patient, userDoctor, doctor, _ := factories.CreateEntries(db)

	// JWT Claims
	claimsDoctor := factories.MakeJWT(userDoctor.ID, models.RoleDoctor)
	claimsPatient := factories.MakeJWT(userPatient.ID, models.RolePatient)

	// Routers & API clients
	routerDoctor := setupPrescriptionRouterWithClaims(claimsDoctor)
	clientDoctor := apiclient.NewTestClient(routerDoctor)

	routerPatient := setupPrescriptionRouterWithClaims(claimsPatient)
	clientPatient := apiclient.NewTestClient(routerPatient)

	t.Run("Create Prescription with Medical Record", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
		prescription := factories.SeedPrescription(db, &record.ID, doctor.ID, patient.ID)

		assert.NotNil(t, prescription.ID)
		assert.Equal(t, patient.ID, prescription.PatientID)
		assert.Equal(t, doctor.ID, prescription.DoctorID)
		assert.NotNil(t, prescription.MedicalRecordID)
	})

	t.Run("Create Prescription without Medical Record", func(t *testing.T) {
		prescription := factories.SeedPrescription(db, nil, doctor.ID, patient.ID)

		assert.NotNil(t, prescription.ID)
		assert.Equal(t, patient.ID, prescription.PatientID)
		assert.Equal(t, doctor.ID, prescription.DoctorID)
		assert.Nil(t, prescription.MedicalRecordID)
	})

	t.Run("Get Prescription by ID", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
		prescription := factories.SeedPrescription(db, &record.ID, doctor.ID, patient.ID)

		res := clientDoctor.Get("/prescriptions/"+prescription.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Get Prescriptions by Patient ID", func(t *testing.T) {
		factories.SeedPrescription(db, nil, doctor.ID, patient.ID)

		res := clientPatient.Get("/prescriptions/patient/"+patient.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)

	})

	t.Run("Update Prescription", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
		prescription := factories.SeedPrescription(db, &record.ID, doctor.ID, patient.ID)

		updateBody := map[string]interface{}{
			"dosage":    "2x daily",
			"issued_at": time.Now().Format(time.RFC3339),
		}
		res := clientDoctor.Put("/prescriptions/"+prescription.ID.String(), updateBody, nil)
		assert.Equal(t, http.StatusOK, res.Code)

	})

	t.Run("Delete Prescription", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
		prescription := factories.SeedPrescription(db, &record.ID, doctor.ID, patient.ID)

		res := clientDoctor.Delete("/prescriptions/"+prescription.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)

	})
}

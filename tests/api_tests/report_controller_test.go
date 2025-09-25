package apitests

import (
	"net/http"
	"testing"

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

func setupReportRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	reportCache := cache.NewCache(config.Rdb, config.Ctx)
	rg := r.Group("/reports")
	routes.RegisterReportRoute(rg, reportCache)
	return r
}

func TestReportRoutes(t *testing.T) {
	db := config.DB

	// Seed users
	userPatient, _, userDoctor, _, _ := factories.CreateEntries(db)

	// JWT Claims
	claimsDoctor := factories.MakeJWT(userDoctor.ID, models.RoleDoctor)
	claimsPatient := factories.MakeJWT(userPatient.ID, models.RolePatient)

	// Routers & API clients
	routerDoctor := setupReportRouterWithClaims(claimsDoctor)
	clientDoctor := apiclient.NewTestClient(routerDoctor)

	routerPatient := setupReportRouterWithClaims(claimsPatient)
	clientPatient := apiclient.NewTestClient(routerPatient)

	t.Run("Create Report", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, userPatient.ID, userDoctor.ID)

		body := map[string]interface{}{
			"title":             "Blood Test",
			"description":       "Routine blood test",
			"file_url":          "http://example.com/report.pdf",
			"patient_id":        userPatient.ID,
			"doctor_id":         userDoctor.ID,
			"medical_record_id": record.ID,
		}

		res := clientDoctor.Post("/reports", body, nil)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), "Report created successfully")
	})

	t.Run("Get Report by Patient ID", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, userPatient.ID, userDoctor.ID)
		factories.SeedReport(db, userPatient.ID, userDoctor.ID, &record.ID)

		res := clientPatient.Get("/reports/patient/"+userPatient.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Blood Test")
	})

	t.Run("Get Report by ID", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, userPatient.ID, userDoctor.ID)
		report := factories.SeedReport(db, userPatient.ID, userDoctor.ID, &record.ID)

		res := clientDoctor.Get("/reports/"+report.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Blood Test")
	})

	t.Run("Update Report", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, userPatient.ID, userDoctor.ID)
		report := factories.SeedReport(db, userPatient.ID, userDoctor.ID, &record.ID)

		updateBody := map[string]interface{}{
			"title":             "Updated Blood Test",
			"description":       "Updated description",
			"file_url":          "http://example.com/updated.pdf",
			"patient_id":        userPatient.ID,
			"doctor_id":         userDoctor.ID,
			"medical_record_id": record.ID,
		}

		res := clientDoctor.Put("/reports/"+report.ID.String(), updateBody, nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Report updated successfully")
	})
	t.Run("Delete Report", func(t *testing.T) {
		record := factories.CreateMedicalRecord(db, userPatient.ID, userDoctor.ID)
		report := factories.SeedReport(db, userPatient.ID, userDoctor.ID, &record.ID)

		res := clientDoctor.Delete("/reports/"+report.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Report deleted successfully")
	})
}

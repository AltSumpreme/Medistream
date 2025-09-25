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

func setupVitalsRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	vitalsCache := cache.NewCache(config.Rdb, config.Ctx)
	rg := r.Group("/vitals")
	routes.RegisterVitalsRoutes(rg, vitalsCache)
	return r
}

func TestVitalRoutes(t *testing.T) {
	db := config.DB

	// Seed users
	userPatient, _, userDoctor, _, _ := factories.CreateEntries(db)

	// JWT Claims
	claimsDoctor := factories.MakeJWT(userDoctor.ID, models.RoleDoctor)
	claimsPatient := factories.MakeJWT(userPatient.ID, models.RolePatient)

	// Routers & API clients
	routerDoctor := setupVitalsRouterWithClaims(claimsDoctor)
	clientDoctor := apiclient.NewTestClient(routerDoctor)

	routerPatient := setupVitalsRouterWithClaims(claimsPatient)
	clientPatient := apiclient.NewTestClient(routerPatient)

	t.Run("Create Vital", func(t *testing.T) {
		body := map[string]interface{}{
			"patient_id":  userPatient.ID,
			"type":        "HEART_RATE",
			"value":       "85",
			"status":      "normal",
			"recorded_at": time.Now().Format(time.RFC3339),
		}
		res := clientDoctor.Post("/vitals/", body, nil)
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), "vital_id")
	})

	t.Run("Get Vital by ID", func(t *testing.T) {
		vital := factories.SeedVital(db, userPatient.ID)

		res := clientDoctor.Get("/vitals/"+vital.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), vital.Type)
	})

	t.Run("Get Vitals by Patient ID", func(t *testing.T) {
		factories.SeedVital(db, userPatient.ID)

		res := clientPatient.Get("/vitals/patient/"+userPatient.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "HEART_RATE")
	})

	t.Run("Update Vital", func(t *testing.T) {
		vital := factories.SeedVital(db, userPatient.ID)

		updateBody := map[string]interface{}{
			"value":       "95",
			"status":      "elevated",
			"recorded_at": time.Now().Format(time.RFC3339),
		}
		res := clientDoctor.Put("/vitals/"+vital.ID.String(), updateBody, nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "updated")
	})

	t.Run("Delete Vital", func(t *testing.T) {
		vital := factories.SeedVital(db, userPatient.ID)

		res := clientDoctor.Delete("/vitals/"+vital.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "deleted")
	})
}

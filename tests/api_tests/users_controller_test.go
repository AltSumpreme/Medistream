package apitests

import (
	"net/http"
	"testing"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/routes"
	apiclient "github.com/AltSumpreme/Medistream.git/tests/api_client"
	"github.com/AltSumpreme/Medistream.git/tests/factories"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	rg := r.Group("/users")
	routes.RegisterUserRoutes(rg)
	return r
}

func TestUserRoutes(t *testing.T) {
	db := config.DB

	// Seed entries
	userPatient, _, _, _, _ := factories.CreateEntries(db)

	claimsPatient := factories.MakeJWT(userPatient.ID, models.RolePatient)
	// claimsDoctor := factories.MakeJWT(userDoctor.ID, models.RoleDoctor)

	routerPatient := setupUserRouterWithClaims(claimsPatient)
	clientPatient := apiclient.NewTestClient(routerPatient)

	//routerDoctor := setupUserRouterWithClaims(claimsDoctor)
	// clientDoctor := apiclient.NewTestClient(routerDoctor)

	t.Run("Get User Profile", func(t *testing.T) {
		res := clientPatient.Get("/users/"+userPatient.ID.String(), nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), userPatient.FirstName)
	})

	t.Run("Update User Profile", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"firstName": "Updated",
			"lastName":  "User",
		}
		res := clientPatient.Put("/users/"+userPatient.ID.String(), updateBody, nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "User profile updated successfully")
	})

	t.Run("Promote Patient to Doctor", func(t *testing.T) {
		newPatient, _, _, _, _ := factories.CreateEntries(db)

		claimsNewPatient := factories.MakeJWT(newPatient.ID, models.RolePatient)
		routerNewPatient := setupUserRouterWithClaims(claimsNewPatient)
		clientNewPatient := apiclient.NewTestClient(routerNewPatient)

		res := clientNewPatient.Put("/users/promote/"+newPatient.ID.String()+"?specialization=DERMATOLOGY", nil, nil)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "User promoted to doctor successfully")
	})

	t.Run("Get Doctors by Specialization", func(t *testing.T) {
		// Ensure doctor exists with specialization
		doctorUser := factories.SeedUser(db, models.RoleDoctor)
		doc := factories.SeedDoctor(db, doctorUser)
		factories.ChangeDoctorSpecialization(db, doc, "CARDIOLOGY")

		res := clientPatient.Get("/users/doctors?specialization=CARDIOLOGY", nil)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "CARDIOLOGY")
	})
}

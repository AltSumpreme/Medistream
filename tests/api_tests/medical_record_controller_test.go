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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupMRRouterWithClaims(claims *utils.JWTClaims) *gin.Engine {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("jwtPayload", claims)
		c.Next()
	})
	rg := r.Group("/medical-records")
	routes.RegisterMedicalRecordsRoutes(rg)
	return r
}

func TestCreateMedicalRecordRoutes(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)
	body := map[string]interface{}{
		"patient_id": patient.ID,
		"doctor_id":  doctor.ID,
		"diagnosis":  "Routine checkup",
		"notes":      "All good",
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Post("/medical-records/", body, headers)
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Contains(t, res.Body.String(), "Medical record created successfully")

}

func TestGetAllMedicalrecords(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Get("/medical-records", headers)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetMedicalRecordByID(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Get("/medical-records/"+record.ID.String(), headers)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetRecordsByPatientID(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Get("/medical-records/patient/"+patient.ID.String(), headers)
	assert.Equal(t, http.StatusOK, res.Code)

}

func TestUpdateMedicalRecords(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)
	body := map[string]interface{}{
		"diagnosis": "Updated Diagnosis",
		"notes":     "Updated Notes",
	}
	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Put("/medical-records/"+record.ID.String(), body, headers)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Medical record updated successfully")
}

func TestSofttDeleteMedicalRecord(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleDoctor)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Delete("/medicalrecords/soft-delete/"+record.ID.String(), headers)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Medical record deleted successfully")
}

func TestHardDeleteMedicalRecord(t *testing.T) {
	db := config.DB
	claims := makeJWT(uuid.New(), models.RoleAdmin)

	_, patient, _, doctor, _ := factories.CreateEntries(db)
	record := factories.CreateMedicalRecord(db, patient.ID, doctor.ID)
	router := setupMRRouterWithClaims(claims)
	client := apiclient.NewTestClient(router)

	headers := map[string]string{"Content-Type": "application/json"}

	res := client.Delete("/medicalrecords/hard-delete/"+record.ID.String(), headers)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Medical record permanently deleted")
}

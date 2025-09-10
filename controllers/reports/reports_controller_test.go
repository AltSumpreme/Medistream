package reports_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/reports"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testDoctor models.User
var testPatient models.User

func SeedUser(role models.Role) models.User {
	email := fmt.Sprintf("test_%s@example.com", role)
	user := models.User{}
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}
	user = models.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  string(role),
		Role:      role,
	}
	config.DB.Create(&user)
	return user
}

func setupReportTestEnv() *gin.Engine {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()
	config.ConnectDB()

	testDoctor = SeedUser(models.RoleDoctor)
	testPatient = SeedUser(models.RolePatient)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user", &testDoctor)
		c.Next()
	})

	r.POST("/reports", reports.CreateReport)
	r.GET("/reports/patient/:patient_id", reports.GetReportByPatientID)
	r.GET("/reports/:id", reports.GetReportByID)
	r.PUT("/reports/:id", reports.UpdateReportByID)
	r.DELETE("/reports/:id", reports.DeleteReportByID)

	return r
}

func createTestReport(t *testing.T, r *gin.Engine) uuid.UUID {
	record := models.MedicalRecord{
		ID:        uuid.New(),
		PatientID: testPatient.ID,
		DoctorID:  testDoctor.ID,
		Notes:     "Sample",
	}
	config.DB.Create(&record)

	input := reports.ReportInput{
		Title:           "X-Ray",
		Description:     "Lung scan",
		FileURL:         "https://example.com/xray.pdf",
		PatientID:       testPatient.ID,
		MedicalRecordID: record.ID,
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/reports", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var created models.Report
	config.DB.Last(&created)
	return created.ID
}

func TestCreateReport(t *testing.T) {
	r := setupReportTestEnv()
	_ = createTestReport(t, r)
}

func TestGetReportByID(t *testing.T) {
	r := setupReportTestEnv()
	reportID := createTestReport(t, r)

	req := httptest.NewRequest("GET", "/reports/"+reportID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "X-Ray")
}

func TestGetReportsByPatientID(t *testing.T) {
	r := setupReportTestEnv()
	_ = createTestReport(t, r)

	req := httptest.NewRequest("GET", "/reports/patient/"+testPatient.ID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "reports")
}

func TestUpdateReportByID(t *testing.T) {
	r := setupReportTestEnv()
	reportID := createTestReport(t, r)

	update := reports.ReportInput{
		Title:           "Updated Title",
		Description:     "Updated Desc",
		FileURL:         "https://example.com/updated.pdf",
		PatientID:       testPatient.ID,
		MedicalRecordID: uuid.New(), // for test only
	}
	config.DB.Create(&models.MedicalRecord{ID: update.MedicalRecordID, DoctorID: testDoctor.ID, PatientID: testPatient.ID, Notes: "upd"})
	b, _ := json.Marshal(update)
	req := httptest.NewRequest("PUT", "/reports/"+reportID.String(), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestDeleteReportByID(t *testing.T) {
	r := setupReportTestEnv()
	reportID := createTestReport(t, r)

	req := httptest.NewRequest("DELETE", "/reports/"+reportID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted")
}

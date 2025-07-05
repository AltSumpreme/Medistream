package vitals_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/vitals"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	TestPatient models.User
	TestDoctor  models.User
	TestAdmin   models.User
)

func init() {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}
	config.ConnectDB()

	TestPatient = seedTestUser(models.RolePatient)
	TestDoctor = seedTestUser(models.RoleDoctor)
	TestAdmin = seedTestUser(models.RoleAdmin)
}

func seedTestUser(role models.Role) models.User {
	email := "vital_test_user@example.com"
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err == nil {
		return user
	}
	user = models.User{
		ID:        uuid.New(),
		FirstName: "Vital",
		LastName:  "Test",
		Email:     email,
		Password:  "test123",
		Role:      role,
	}
	config.DB.Create(&user)
	config.DB.FirstOrCreate(&models.Patient{ID: user.ID, UserID: user.ID})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		user, ok := val.(*models.User)
		if !ok || user == nil || !slices.Contains(allowedRoles, user.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

func setupRouter(user models.User) *gin.Engine {
	r := gin.Default()
	r.Use(injectTestUser(user))

	v := r.Group("/vitals")
	v.POST("/", testRoleChecker(models.RoleDoctor, models.RoleAdmin), vitals.CreateVital)
	v.GET("/patient/:id", testRoleChecker(models.RolePatient, models.RoleDoctor, models.RoleAdmin), vitals.GetVitalsByPatientID)
	v.GET("/:id", testRoleChecker(models.RolePatient, models.RoleDoctor, models.RoleAdmin), vitals.GetVitalByID)
	v.PUT("/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), vitals.UpdateVital)
	v.DELETE("/:id", testRoleChecker(models.RoleDoctor, models.RoleAdmin), vitals.DeleteVital)
	return r
}

func createTestVital(t *testing.T, router *gin.Engine, patientID uuid.UUID) uuid.UUID {
	v := map[string]interface{}{
		"patient_id":  patientID,
		"type":        "HEART_RATE",
		"value":       "85",
		"status":      "normal",
		"recorded_at": time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(v)

	req := httptest.NewRequest("POST", "/vitals/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)

	idStr := res["vital_id"].(string)
	id, err := uuid.Parse(idStr)
	assert.NoError(t, err)
	return id
}

// ------------------ TESTS -----------------------

func TestCreateVital(t *testing.T) {
	router := setupRouter(TestPatient)
	createTestVital(t, router, TestPatient.ID)
}

func TestGetVitalByID(t *testing.T) {
	router := setupRouter(TestPatient)
	vitalID := createTestVital(t, router, TestPatient.ID)

	req := httptest.NewRequest("GET", "/vitals/"+vitalID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "HEART_RATE")
}

func TestGetVitalsByPatientID(t *testing.T) {
	router := setupRouter(TestPatient)
	_ = createTestVital(t, router, TestPatient.ID)

	req := httptest.NewRequest("GET", "/vitals/patient/"+TestPatient.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "HEART_RATE")
}

func TestUpdateVital(t *testing.T) {
	router := setupRouter(TestPatient)
	vitalID := createTestVital(t, router, TestPatient.ID)

	update := map[string]interface{}{
		"value":       "95",
		"status":      "elevated",
		"recorded_at": time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/vitals/"+vitalID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestDeleteVital(t *testing.T) {
	router := setupRouter(TestPatient)
	vitalID := createTestVital(t, router, TestPatient.ID)

	req := httptest.NewRequest("DELETE", "/vitals/"+vitalID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "soft deleted")
}

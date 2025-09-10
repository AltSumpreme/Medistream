package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/controllers/user"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	utils.InitLogger()
	config.ConnectDB()
}

func setupUserRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users/:id", user.GetUserProfile)
	r.PUT("/users/:id", user.UpdateUserProfile)
	r.POST("/users/promote/:id", user.PromotePatienttoDoctor)
	return r
}

// helper: create an Auth record
func createTestAuth(email string) models.Auth {
	auth := models.Auth{
		ID:       uuid.New(),
		Email:    email,
		Password: "hashedpassword",
	}
	if err := config.DB.Create(&auth).Error; err != nil {
		panic(err)
	}
	return auth
}

// helper: create a User with Auth
func createTestUser(role models.Role) models.User {
	auth := createTestAuth(fmt.Sprintf("user+%d@example.com", time.Now().UnixNano()))

	user := models.User{
		ID:        uuid.New(),
		AuthID:    auth.ID,
		FirstName: "Test",
		LastName:  "User",
		Role:      role,
		Phone:     "1234567890",
	}
	if err := config.DB.Create(&user).Error; err != nil {
		panic(err)
	}
	return user
}

func TestGetUserProfile(t *testing.T) {
	router := setupUserRouter()
	user := createTestUser(models.RolePatient)

	req := httptest.NewRequest("GET", "/users/"+user.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	var body map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &body)
	assert.NoError(t, err)
}

func TestUpdateUserProfile(t *testing.T) {
	router := setupUserRouter()
	user := createTestUser(models.RolePatient)

	update := map[string]string{
		"firstName": "Updated",
		"lastName":  "User",
		"email":     fmt.Sprintf("updated+%d@example.com", time.Now().UnixNano()),
	}
	jsonBody, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/users/"+user.ID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	// Fetch updated user from DB
	var updated models.User
	err := config.DB.First(&updated, "id = ?", user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.FirstName)
}

func TestPromotePatienttoDoctor(t *testing.T) {
	// Create a patient user
	patient := createTestUser(models.RolePatient)

	// Create a test router
	r := gin.Default()
	r.PUT("/promote/:id", user.PromotePatienttoDoctor)

	req, _ := http.NewRequest(http.MethodPut, "/promote/"+patient.ID.String()+"?specialization=GENERAL", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User promoted to doctor successfully")

	// Verify DB changes
	var updatedUser models.User
	err := config.DB.First(&updatedUser, "id = ?", patient.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, models.RoleDoctor, updatedUser.Role)

	var doctor models.Doctor
	err = config.DB.First(&doctor, "user_id = ?", patient.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "GENERAL", doctor.Specialization)
}

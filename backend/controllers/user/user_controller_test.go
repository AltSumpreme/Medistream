package user

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
	r.GET("/users/:id", GetUserProfile)
	r.PUT("/users/:id", UpdateUserProfile)
	r.POST("/users/promote/:id", PromoteUser)
	return r
}

func createTestUser(role models.Role) models.User {
	tx := config.DB.Begin()

	defer tx.Rollback()
	user := models.User{
		ID:        uuid.New(),
		FirstName: "Test",
		LastName:  "User",
		Role:      role,
	}
	config.DB.Create(&user)
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

	tx := config.DB.Begin()
	defer tx.Rollback()
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
	_ = config.DB.First(&updated, "id = ?", user.ID)
	assert.Equal(t, "Updated", updated.FirstName)
}

func TestPromoteUser(t *testing.T) {
	router := setupUserRouter()

	tx := config.DB.Begin()
	defer tx.Rollback()

	// Patient -> Doctor
	patient := createTestUser(models.RolePatient)
	assert.Equal(t, models.RolePatient, patient.Role)
	req := httptest.NewRequest("POST", "/users/promote/"+patient.ID.String(), nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	var promoted models.User
	config.DB.First(&promoted, "id = ?", patient.ID)
	assert.Equal(t, models.RoleDoctor, promoted.Role)

	// Doctor -> Admin
	{ /*	req = httptest.NewRequest("POST", "/users/promote/"+patient.ID.String(), nil)
			res = httptest.NewRecorder()
			router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Code)

			config.DB.First(&promoted, "id = ?", patient.ID)
			assert.Equal(t, models.RoleAdmin, promoted.Role) */
	}

	// Admin -> Error
	{ /*req = httptest.NewRequest("POST", "/users/promote/"+patient.ID.String(), nil)
		res = httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code) */
	}
}

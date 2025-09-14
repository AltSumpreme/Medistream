package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	utils.InitLogger()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	config.ConnectDB()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/signup", SignUp)
	r.POST("/login", Login)

	r.POST("/refresh", RefreshAccessToken)
	r.POST("/logout", Logout)
	return r
}

func TestAuthFlow(t *testing.T) {
	router := setupRouter()

	email := "john.doe+" + time.Now().Format("150405") + "@example.com"
	password := "securePassword123!"

	var accessToken, refreshToken string

	t.Run("Signup", func(t *testing.T) {
		signUp := map[string]string{
			"firstname": "John",
			"lastname":  "Doe",
			"email":     email,
			"password":  password,
			"phone":     "1234567890",
		}
		body, _ := json.Marshal(signUp)

		req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code, "Signup failed")
	})

	t.Run("Login", func(t *testing.T) {
		login := map[string]string{
			"email":    email,
			"password": password,
		}
		body, _ := json.Marshal(login)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code, "Login failed")

		var loginResp map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &loginResp)
		assert.NoError(t, err, "Failed to parse login response")

		at, ok := loginResp["access_token"].(string)
		assert.True(t, ok && at != "", "Access token missing")
		accessToken = at

		// manually get refresh token from DB
		var user models.User
		err = config.DB.Where("email = ?", email).First(&user).Error
		assert.NoError(t, err, "User not found after signup")

		var refresh models.RefreshToken
		err = config.DB.Where("user_id = ?", user.ID).Last(&refresh).Error
		assert.NoError(t, err, "Refresh token not found")

		refreshToken = refresh.Token
	})

	t.Run("Verify Access Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/verify", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code, "Token verification failed")
	})

	t.Run("Refresh Access Token", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"refresh_token": refreshToken,
		})

		req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code, "Access token refresh failed")
	})

	t.Run("Logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code, "Logout failed")
	})
}

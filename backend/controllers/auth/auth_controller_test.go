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
	r.GET("/verify", VerifyToken)
	r.POST("/refresh", RefreshAccessToken)
	r.POST("/logout", Logout)
	return r
}

func TestAuthFlow(t *testing.T) {
	router := setupRouter()

	// unique email
	email := "john.doe+" + time.Now().Format("150405") + "@example.com"
	password := "securePassword123!"

	// ---------------- SIGNUP ----------------
	signUp := map[string]string{
		"firstname": "John",
		"lastname":  "Doe",
		"email":     email,
		"password":  password,
	}
	signUpBody, _ := json.Marshal(signUp)
	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(signUpBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	// ---------------- LOGIN ----------------
	login := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBody, _ := json.Marshal(login)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var loginResp map[string]interface{}
	_ = json.Unmarshal(res.Body.Bytes(), &loginResp)

	accessToken := loginResp["access_token"].(string)
	assert.NotEmpty(t, accessToken)

	// manually fetch refresh token from DB for verification
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	assert.NoError(t, err)

	var refresh models.RefreshToken
	err = config.DB.Where("user_id = ?", user.ID).Last(&refresh).Error
	assert.NoError(t, err)
	refreshToken := refresh.Token

	// ---------------- VERIFY TOKEN ----------------
	req = httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	// ---------------- REFRESH ACCESS TOKEN ----------------
	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}
	refreshBody, _ := json.Marshal(refreshReq)
	req = httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	// ---------------- LOGOUT ----------------
	req = httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

package apitests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/routes"
	apiclient "github.com/AltSumpreme/Medistream.git/tests/api_client"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter() *gin.Engine {
	r := gin.Default()
	rg := r.Group("/auth")
	routes.RegisterAuthRoutes(rg)
	return r
}

func TestAuthFlow(t *testing.T) {
	router := setupAuthRouter()
	client := apiclient.NewTestClient(router)

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

		res := client.Post("/auth/signup", signUp, nil)
		assert.Equal(t, http.StatusOK, res.Code, "Signup failed")
	})

	t.Run("Login", func(t *testing.T) {
		login := map[string]string{
			"email":    email,
			"password": password,
		}

		res := client.Post("/auth/login", login, nil)
		assert.Equal(t, http.StatusOK, res.Code, "Login failed")

		var loginResp map[string]interface{}
		err := json.Unmarshal(res.Body.Bytes(), &loginResp)
		assert.NoError(t, err, "Failed to parse login response")

		at, ok := loginResp["access_token"].(string)
		assert.True(t, ok && at != "", "Access token missing")
		accessToken = at

		// Manually get refresh token from DB
		var auth models.Auth
		err = config.DB.Where("email = ?", email).First(&auth).Error
		assert.NoError(t, err, "User not found after signup")

		var refresh models.RefreshToken
		var user models.User
		err = config.DB.Where("auth_id = ?", auth.ID).First(&user).Error
		assert.NoError(t, err, "User not found after signup")
		err = config.DB.Where("user_id = ?", user.ID).Last(&refresh).Error
		assert.NoError(t, err, "Refresh token not found")

		refreshToken = refresh.Token
	})

	t.Run("Refresh Access Token", func(t *testing.T) {
		body := map[string]string{
			"refresh_token": refreshToken,
		}
		res := client.Post("/auth/refresh", body, nil)
		assert.Equal(t, http.StatusOK, res.Code, "Access token refresh failed")
	})

	t.Run("Logout", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		res := client.Post("/auth/logout", nil, headers)
		assert.Equal(t, http.StatusOK, res.Code, "Logout failed")
	})
}

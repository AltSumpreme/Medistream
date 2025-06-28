package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	log.Println("SignUp called")
	var input struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string ` json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {

		utils.Log.Warnf("SignUp: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.Log.Warnf("SignUp: Email already exists - %s", input.Email)
		// If the user already exists, return an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedpassword, err := utils.HashPassword(input.Password)

	if err != nil {
		utils.Log.Warnf("Signup: Hash Password - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  string(hashedpassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("Signup: Failed to create user %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	patient := models.Patient{
		UserID: user.ID,
	}
	if err := config.DB.Create(&patient).Error; err != nil {
		utils.Log.Errorf("SignUp: Failed to create patient profile - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient profile"})
		return
	}
	utils.Log.Infof("SignUp: User succesfully signed up")
	c.JSON(http.StatusOK, gin.H{"message": "User signed up successfully"})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Login:Invalid input %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Log.Errorf("Login:Invalid email or password- %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, input.Password); err != nil {
		utils.Log.Warnf("Login:Invalid password - %v,", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid  password"})
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		utils.Log.Errorf("Login: Failed to generate token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.Log.Errorf("Login: Failed to generate refresh token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	rt := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := config.DB.Create(&rt).Error; err != nil {
		utils.Log.Errorf("Login: Failed to create refresh token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "User logged in successfully",
		"access_token": accessToken,
	})
}

func VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Log.Warnf("Verify Token: Missing or malformed token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		utils.Log.Warnf("Verify Token: Failed to validate JWT Token - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "token is valid",
		"user_id":  claims.UserID,
		"role":     claims.Role,
		"issuedAt": claims.IssuedAt,
		"expires":  claims.ExpiresAt,
	})
}

func RefreshAccessToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Refresh Access Token: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stored models.RefreshToken
	if err := config.DB.Preload("User").Where("token = ? AND revoked = false", input.RefreshToken).First(&stored).Error; err != nil {
		utils.Log.Errorf("RefreshAccessToken:Invalid or expired refresh token - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	if time.Now().After(stored.ExpiresAt) {
		utils.Log.Warnf("RefreshAccessToken:Refresh token has expired")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
		return
	}

	accessToken, _ := utils.GenerateJWT(stored.UserID, string(stored.User.Role))
	c.JSON(http.StatusOK, gin.H{
		"message":      "Access token refreshed successfully",
		"access_token": accessToken})

}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	var stored models.RefreshToken
	if err := config.DB.Where("token = ? AND revoked = false", tokenStr).First(&stored).Error; err != nil {
		utils.Log.Errorf("Logout:Invalid or expired refresh token-%v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	stored.Revoked = true
	if err := config.DB.Save(&stored).Error; err != nil {
		utils.Log.Errorf("Logout:Failed to revoke token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

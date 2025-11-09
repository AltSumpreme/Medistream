package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/metrics"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/queue"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	var input struct {
		FirstName string `json:"firstname" binding:"required"`
		LastName  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		Phone     string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("SignUp: Invalid input - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAuth models.Auth
	err := metrics.DbMetrics(config.DB, "Signup", func(db *gorm.DB) error {
		return db.Where("email = ?", input.Email).First(&existingAuth).Error
	})
	if err == nil {
		utils.Log.Warnf("SignUp: Email already exists - %s", input.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.Log.Warnf("Signup: Hash Password - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create Auth
	auth := models.Auth{
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := metrics.DbMetrics(config.DB, "Signup", func(db *gorm.DB) error {
		return db.Create(&auth).Error
	}); err != nil {
		utils.Log.Errorf("Signup: Failed to create auth %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create User
	user := models.User{
		AuthID:    auth.ID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      models.RolePatient,
		Phone:     input.Phone,
	}

	if err := metrics.DbMetrics(config.DB, "Signup", func(db *gorm.DB) error {
		return db.Create(&user).Error
	}); err != nil {
		utils.Log.Errorf("SignUp: Failed to create user profile - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	// Create Patient record
	patient := models.Patient{UserID: user.ID}
	if err := metrics.DbMetrics(config.DB, "Signup", func(db *gorm.DB) error {
		return db.Create(&patient).Error
	}); err != nil {
		utils.Log.Errorf("SignUp: Failed to create patient profile - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient profile"})
		return
	}

	task, err := queue.NewWelcomeEmailTask(input.Email, input.FirstName)
	if err != nil {
		utils.Log.Errorf("SignUp: Failed to create welcome email task - %v", err)
	} else {
		_, enqueueErr := queue.Client.Enqueue(
			task,
			asynq.Queue("emails"),
			asynq.MaxRetry(3),
		)

		if enqueueErr != nil {
			utils.Log.Errorf("SignUp: Failed to enqueue welcome email - %v", enqueueErr)
		}
	}

	utils.Log.Infof("SignUp: User successfully signed up")
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
	var auth models.Auth
	err := metrics.DbMetrics(config.DB, "Login", func(db *gorm.DB) error {
		return db.Where("email = ?", input.Email).First(&auth).Error
	})
	if err != nil {
		utils.Log.Errorf("Login:Invalid email or password- %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := utils.VerifyPassword(auth.Password, input.Password); err != nil {
		utils.Log.Warnf("Login:Invalid password - %v,", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid  password"})
		return
	}

	var user models.User
	if err := metrics.DbMetrics(config.DB, "Login", func(db *gorm.DB) error {
		return db.Where("auth_id = ?", auth.ID).First(&user).Error
	}); err != nil {
		utils.Log.Errorf("Login:Failed to find user profile - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user profile"})
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
	err = metrics.DbMetrics(config.DB, "Login", func(db *gorm.DB) error {
		return db.Create(&rt).Error
	})
	if err != nil {
		utils.Log.Errorf("Login: Failed to create refresh token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "User logged in successfully",
		"access_token": accessToken,
	})
}

/*
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
			"token":    tokenStr,
		})
	}
*/
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
	err := metrics.DbMetrics(config.DB, "RefreshAccessToken", func(db *gorm.DB) error {
		return db.Preload("User").Where("token = ? AND revoked = false", input.RefreshToken).First(&stored).Error
	})
	if err != nil {
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
	c.SetCookie("access_token", accessToken, 7200, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Access token refreshed successfully",
	})

}

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		utils.Log.Warnf("Logout:Invalid or Expired token")
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired  token"})
		return
	}
	err = metrics.DbMetrics(config.DB, "Logout", func(db *gorm.DB) error {
		return db.Model(&models.RefreshToken{}).Where("user_id=? AND revoked = false", claims.UserID).Update("revoked", true).Error
	})
	if err != nil {
		utils.Log.Errorf("Logout: Failed to revoke refresh token - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh token"})
		return
	}
	c.SetCookie("access_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var auth models.Auth
	err := metrics.DbMetrics(config.DB, "ForgotPassword", func(db *gorm.DB) error {
		return db.Where("email = ?", input.Email).First(&auth).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		utils.Log.Errorf("ForgotPassword: Failed to generate OTP - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}
	if err := cache.SaveOTP(input.Email, otp, 10); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	task, err := queue.NewOTPEmailTask(input.Email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	_, err = queue.Client.Enqueue(task, asynq.Queue("emails"), asynq.MaxRetry(3))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email"})
}

func ResendOTP(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var auth models.Auth
	err := metrics.DbMetrics(config.DB, "ResendOTP", func(db *gorm.DB) error {
		return db.Where("email = ?", input.Email).First(&auth).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
		return
	}

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		utils.Log.Errorf("ResendOTP: Failed to generate OTP - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}
	if err := cache.SaveOTP(input.Email, otp, 10); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
		return
	}

	task, err := queue.NewOTPEmailTask(input.Email, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email task"})
		return
	}

	_, err = queue.Client.Enqueue(task, asynq.Queue("emails"), asynq.MaxRetry(3))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}

func ResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		ResetToken  string `json:"resetToken" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user from DB
	var auth models.Auth
	if err := config.DB.Where("email = ?", req.Email).First(&auth).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or token"})
		return
	}

	// Verify OTP token from Redis
	if err := cache.VerifyOTP(req.Email, req.ResetToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update user password in DB
	auth.Password = hashedPassword
	if err := config.DB.Save(&auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Push password reset confirmation email to queue
	emailTemplate := utils.GetPasswordResetSuccessTemplate()
	task, err := queue.ResetEmailTask(auth.Email, emailTemplate.Subject, emailTemplate.Body)
	if err != nil {
		utils.Log.Warnf("ResetPassword: failed to create email task for %s: %v", auth.Email, err)
	}
	if _, err := queue.Client.Enqueue(task, asynq.Queue("emails"), asynq.MaxRetry(3)); err != nil {
		utils.Log.Warnf("ResetPassword: failed to enqueue email for %s: %v", auth.Email, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully",
	})
}

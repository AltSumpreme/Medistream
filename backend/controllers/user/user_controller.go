package user

import (
	"net/http"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	var userID = c.Param("id")

	user := models.User{}

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{
		"id":        user.ID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}
func UpdateUserProfile(c *gin.Context) {
	var userID = c.Param("id")
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User profile updated successfully",
	})
	c.JSON(200, gin.H{
		"id":        user.ID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})

}
func PromoteUser(c *gin.Context) {
	var userID = c.Param("id")

	var user models.User

	if err := config.DB.Where("id=?", userID).First(&user).Error; err != nil {
		utils.Log.Warnf("PromoteUser: User not found - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	switch user.Role {
	case models.RolePatient:
		user.Role = models.RoleDoctor

	case models.RoleDoctor:

		user.Role = models.RoleAdmin
	case models.RoleAdmin:
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is aldready an admin"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
		return
	}
	if err := config.DB.Save(&user).Error; err != nil {
		utils.Log.Errorf("PromoteUser: Failed to promote user - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User promoted successfully"})
}

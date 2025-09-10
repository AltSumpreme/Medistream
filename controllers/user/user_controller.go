package user

import (
	"net/http"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

func GetDoctorsBySpecialization(c *gin.Context) {
	specialization := c.Query("specialization")
	if specialization == "" {
		utils.Log.Warn("GetDoctorsBySpecialization: Specialization Query is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specialization is required"})
		return
	}

	var doctors []models.Doctor
	if err := config.DB.Preload("User").Where("specialization = ?", specialization).Find(&doctors).Error; err != nil {
		utils.Log.Errorf("GetDoctorsBySpecialization: Failed to fetch doctors - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type DoctorResponse struct {
		DoctorID uuid.UUID `json:"doctor_id"`
		Name     string    `json:"name"`
	}

	var response []DoctorResponse
	for _, doc := range doctors {
		if doc.User != nil {
			response = append(response, DoctorResponse{
				DoctorID: doc.ID,
				Name:     doc.User.FirstName + " " + doc.User.LastName,
			})
		}
	}

	c.JSON(http.StatusOK, response)
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

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user profile"})
		return
	}
	response := gin.H{
		"message":    "User profile updated successfully",
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	if user.Role == models.RoleDoctor {
		if input.Doctor != nil {
			var doctor models.Doctor
			if err := config.DB.First(&doctor, "user_id = ?", user.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor profile not found"})
				return
			}

			doctor.Specialization = input.Doctor.Specialization
			if err := config.DB.Save(&doctor).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update specialization"})
				return
			}

			response["specialization"] = doctor.Specialization
		}
	}

	c.JSON(200, gin.H{
		"message": "User profile updated successfully",
	})
	c.JSON(200, gin.H{
		"id":        user.ID,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})

}

func PromotePatienttoDoctor(c *gin.Context) {
	var UserID = c.Param("id")
	var specialization = c.Query("specialization")
	if specialization == "" {
		utils.Log.Warn("PromotePatienttoDoctor: Specialization parameter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Specialization is required"})
		return
	}

	var user models.User

	if err := config.DB.Where("id=?", UserID).First(&user).Error; err != nil {
		utils.Log.Warnf("PromotePatienttoDoctor: User not found - %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if user.Role != models.RolePatient {
		utils.Log.Warnf("PromotePatienttoDoctor: User is not a patient - %v", user.Role)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a patient"})
		return
	}
	user.Role = models.RoleDoctor

	if err := config.DB.Save(&user).Error; err != nil {
		utils.Log.Errorf("PromotePatienttoDoctor: Failed to promote user - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}
	doctor := models.Doctor{
		ID:             uuid.New(),
		UserID:         user.ID,
		Specialization: specialization,
	}
	if err := config.DB.Create(&doctor).Error; err != nil {
		utils.Log.Errorf("PromotePatienttoDoctor: Failed to create doctor profile - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create doctor profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to doctor successfully"})
}

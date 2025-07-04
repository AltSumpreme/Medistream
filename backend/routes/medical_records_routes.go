package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/medicalrecords"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterMedicalRecordsRoutes(router *gin.Engine) {
	medicalRecordsGroup := router.Group("/api/medical-records")
	{
		medicalRecordsGroup.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.CreateMedicalRecord)
		medicalRecordsGroup.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.GetMedicalRecord)
		medicalRecordsGroup.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.UpdateMedicalRecord)
		medicalRecordsGroup.DELETE("/soft/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.SoftDeleteMedicalRecord)
		medicalRecordsGroup.DELETE("/hard/:id", utils.RoleChecker(models.RoleAdmin), medicalrecords.HardDeleteMedicalRecord)
	}
}

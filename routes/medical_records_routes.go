package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/medicalrecords"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterMedicalRecordsRoutes(rg *gin.RouterGroup) {

	{
		rg.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.CreateMedicalRecord)
		rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.GetMedicalRecord)
		rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.UpdateMedicalRecord)
		rg.DELETE("/soft/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.SoftDeleteMedicalRecord)
		rg.DELETE("/hard/:id", utils.RoleChecker(models.RoleAdmin), medicalrecords.HardDeleteMedicalRecord)
	}
}

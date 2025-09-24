package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/medicalrecords"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterMedicalRecordsRoutes(rg *gin.RouterGroup, medicalrecordCache *cache.Cache) {

	{
		rg.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.CreateMedicalRecord)
		rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), medicalrecords.GetMedicalRecordByID)
		rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) {
			medicalrecords.UpdateMedicalRecord(c, medicalrecordCache)
		})
		rg.DELETE("/soft-delete/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) {
			medicalrecords.SoftDeleteMedicalRecord(c, medicalrecordCache)
		})
		rg.DELETE("/hard-delete/:id", utils.RoleChecker(models.RoleAdmin), func(c *gin.Context) {
			medicalrecords.HardDeleteMedicalRecord(c, medicalrecordCache)
		})
	}
}

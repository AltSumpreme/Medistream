package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/vitals"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterVitalsRoutes(rg *gin.RouterGroup, vitalsCache *cache.Cache) {
	rg.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), vitals.CreateVital)
	rg.GET("/patient/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), vitals.GetVitalsByPatientID)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), vitals.GetVitalByID)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) {
		vitals.UpdateVital(c, vitalsCache)
	})
	rg.DELETE("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) {
		vitals.DeleteVital(c, vitalsCache)
	})
}

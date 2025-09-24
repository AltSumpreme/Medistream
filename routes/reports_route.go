package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/reports"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterReportRoute(rg *gin.RouterGroup, reportsCache *cache.Cache) {

	rg.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), reports.CreateReport)
	rg.GET("/patient/:patient_id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), reports.GetReportByPatientID)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), reports.GetReportByID)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) { reports.UpdateReportByID(c, reportsCache) })
	rg.DELETE("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), func(c *gin.Context) { reports.DeleteReportByID(c, reportsCache) })

}

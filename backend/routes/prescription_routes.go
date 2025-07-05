package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/prescriptions"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterPrescriptionRoutes(rg *gin.RouterGroup) {
	rg.POST("/", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), prescriptions.CreatePrescription)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), prescriptions.GetPrescriptionByID)
	rg.GET("/patient/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), prescriptions.GetPrescriptionsByPatientID)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), prescriptions.UpdatePrescription)
	rg.DELETE("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor), prescriptions.DeletePrescription)
}

package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(rg *gin.RouterGroup) {

	rg.POST("/", utils.RoleChecker(models.RolePatient, models.RoleDoctor, models.RoleAdmin), controllers.CreateAppointment)
	rg.GET("/", utils.RoleChecker(models.RoleAdmin), controllers.GetAllAppointments)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), controllers.GetAppointmentByID)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), controllers.UpdateAppointment)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin), controllers.ChangeAppointmentStatus)
	rg.PUT("/reschedule/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.RescheduleAppointment)
	rg.PUT("/cancel/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.CancelAppointment)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.GetAppointmentByDoctorID)
	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.GetAppointmentByPatientID)
	rg.DELETE("/:id", utils.RoleChecker(models.RoleAdmin), controllers.DeleteAppointment)

}

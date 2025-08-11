package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/appointments"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(rg *gin.RouterGroup) {

	rg.POST("", utils.RoleChecker(models.RolePatient, models.RoleDoctor, models.RoleAdmin), appointments.CreateAppointment)
	rg.GET("", utils.RoleChecker(models.RoleAdmin), appointments.GetAllAppointments)
	rg.GET(":id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), appointments.GetAppointmentByID)
	rg.PUT(":id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), appointments.UpdateAppointment)
	rg.PUT("status/:id", utils.RoleChecker(models.RoleAdmin), appointments.ChangeAppointmentStatus)
	rg.PUT("reschedule/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.RescheduleAppointment)
	rg.PUT("cancel/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.CancelAppointment)
	rg.GET("doctor/:id", utils.RoleChecker(models.RoleAdmin), appointments.GetAppointmentByDoctorID)
	rg.GET("patient/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.GetAppointmentByPatientID)
	rg.DELETE(":id", utils.RoleChecker(models.RoleAdmin), appointments.DeleteAppointment)

}

package routes

import (
	"github.com/AltSumpreme/Medistream.git/handlers"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func RegisterAppointmentRoutes(rg *gin.RouterGroup, appointmentCache *cache.Cache, queue *asynq.Client) {

	rg.POST("", utils.RoleChecker(models.RolePatient), func(c *gin.Context) { handlers.HandleUserCreateAppointment(c, queue) })
	{ /*
			rg.GET("", utils.RoleChecker(models.RoleAdmin), appointments.GetAllAppointments)
			rg.GET(":id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), appointments.GetAppointmentByID)
			rg.PUT(":id", utils.RoleChecker(models.RoleAdmin, models.RolePatient, models.RoleDoctor), func(c *gin.Context) {
				appointments.UpdateAppointment(c, appointmentCache)
			})
			rg.PUT("status/:id", utils.RoleChecker(models.RoleAdmin), appointments.ChangeAppointmentStatus)
			rg.PUT("reschedule/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), func(c *gin.Context) {
				appointments.RescheduleAppointment(c, appointmentCache)
			})
			rg.PUT("cancel/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.CancelAppointment)
			rg.GET("doctor/:id", utils.RoleChecker(models.RoleAdmin), appointments.GetAppointmentByDoctorID)
			rg.GET("patient/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), appointments.GetAppointmentByPatientID)
			rg.DELETE(":id", utils.RoleChecker(models.RoleAdmin), func(c *gin.Context) {
				appointments.DeleteAppointment(c, appointmentCache)
			})
		*/
	}
}

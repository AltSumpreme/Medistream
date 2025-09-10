package routes

import (
	"github.com/AltSumpreme/Medistream.git/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	RegisterAuthRoutes(auth)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	RegisterUserRoutes(protected.Group("/user"))
	RegisterAppointmentRoutes(protected.Group("/appointments"))
	RegisterMedicalRecordsRoutes(protected.Group("/medical-records"))
	RegisterVitalsRoutes(protected.Group("/vitals"))
	RegisterPrescriptionRoutes(protected.Group("/prescriptions"))
}

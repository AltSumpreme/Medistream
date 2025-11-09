package routes

import (
	"github.com/AltSumpreme/Medistream.git/middleware"
	"github.com/AltSumpreme/Medistream.git/services/cache"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func RegisterRoutes(r *gin.Engine, appointmentCache *cache.Cache, medicalrecordsCache *cache.Cache, prescriptionsCache *cache.Cache, reportsCache *cache.Cache, vitalsCache *cache.Cache, jobQueue *asynq.Client) {

	auth := r.Group("/auth")
	// auth.Use(middleware.StrictRateLimiterMiddleware())
	RegisterAuthRoutes(auth)

	protected := r.Group("/")
	// protected.Use(middleware.RateLimiterMiddleware())
	protected.Use(middleware.AuthMiddleware())

	RegisterUserRoutes(protected.Group("/user"))
	RegisterAppointmentRoutes(protected.Group("/appointments"), appointmentCache, jobQueue)
	RegisterMedicalRecordsRoutes(protected.Group("/medical-records"), medicalrecordsCache)
	RegisterReportRoute(protected.Group("/reports"), reportsCache)
	RegisterVitalsRoutes(protected.Group("/vitals"), vitalsCache)
	RegisterPrescriptionRoutes(protected.Group("/prescriptions"), prescriptionsCache)
}

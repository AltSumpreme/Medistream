package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/signup", controllers.SignUp)
		api.POST("/login", controllers.Login)

	}
}

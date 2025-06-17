package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
		auth.POST("/verify", controllers.VerifyToken)

	}
	user := r.Group("/user")
	{
		user.GET("/:id", controllers.GetUserProfile)
		user.PUT("/:id", controllers.UpdateUserProfile)

	}
}

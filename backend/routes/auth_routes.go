package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {

	{
		rg.POST("/signup", controllers.SignUp)
		rg.POST("/login", controllers.Login)
		rg.POST("/verify", controllers.VerifyToken)
		rg.POST("/refresh", controllers.RefreshAccessToken)
		rg.POST("/logout", controllers.Logout)

	}
}

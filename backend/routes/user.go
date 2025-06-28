package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {

	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.GetUserProfile)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), controllers.UpdateUserProfile)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin), controllers.PromoteUser)

}

package routes

import (
	"github.com/AltSumpreme/Medistream.git/controllers/user"
	"github.com/AltSumpreme/Medistream.git/models"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {

	rg.GET("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), user.GetUserProfile)
	rg.PUT("/:id", utils.RoleChecker(models.RoleAdmin, models.RoleDoctor, models.RolePatient), user.UpdateUserProfile)
	rg.PUT("/promote/:id", utils.RoleChecker(models.RoleAdmin), user.PromoteUser)

}

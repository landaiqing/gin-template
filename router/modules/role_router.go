package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var roleApi = controller.Controller.RoleController

func RoleRouter(router *gin.RouterGroup) {
	group := router.Group("/auth")
	{
		group.POST("/role/create", roleApi.CreateRole)
		group.POST("/role/add_role_to_user", roleApi.AddRoleToUser)
	}

}

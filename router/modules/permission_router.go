package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var permissionApi = controller.Controller.PermissionController

func PermissionRouter(router *gin.RouterGroup) {
	group := router.Group("/auth/permission")
	{
		group.POST("/add", permissionApi.AddPermissions)
		group.POST("/get_user_permissions", permissionApi.GetUserPermissions)
	}
}

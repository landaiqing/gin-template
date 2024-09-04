package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var permissionApi = api.Api.PermissionApi

func PermissionRouter(router *gin.RouterGroup) {
	group := router.Group("/auth/permission")
	{
		group.POST("/add", permissionApi.AddPermissions)
		group.GET("/get_user_permissions", permissionApi.GetUserPermissions)
	}
}

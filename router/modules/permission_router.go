package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var permissionApi = api.Api.PermissionApi

func PermissionRouter(router *gin.RouterGroup) {
	group := router.Group("/auth/permission")
	//group.Use(middleware.JWTAuthMiddleware())
	group.POST("/add", permissionApi.AddPermissions)
}

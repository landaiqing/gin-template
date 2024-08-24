package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
	"schisandra-cloud-album/middleware"
)

var roleApi = api.Api.RoleApi

func RoleRouter(router *gin.RouterGroup) {
	group := router.Group("/auth")
	group.Use(middleware.JWTAuthMiddleware())
	group.POST("/role/create", roleApi.CreateRole)
	group.POST("/role/add_role_to_user", roleApi.AddRoleToUser)
}

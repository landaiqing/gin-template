package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var websocketAPI = api.Api.WebsocketApi

func WebsocketRouter(router *gin.RouterGroup) {
	group := router.Group("/ws")
	{
		group.GET("/socket", websocketAPI.NewSocketClient)
		group.GET("/delete", websocketAPI.DeleteClient)
		group.GET("/gws", websocketAPI.NewGWSServer)
	}

}

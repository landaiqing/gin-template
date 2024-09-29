package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var websocketAPI = controller.Controller.WebsocketController

func WebsocketRouter(router *gin.RouterGroup) {
	group := router.Group("/ws")
	{
		group.GET("/gws", websocketAPI.NewGWSServer)
	}

}

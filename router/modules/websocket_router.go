package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var qrWebsocketAPI = controller.Controller.QrWebsocketController
var messageWebsocketAPI = controller.Controller.MessageWebsocketController

func WebsocketRouter(router *gin.RouterGroup) {
	group := router.Group("/ws")
	{
		group.GET("/qr_ws", qrWebsocketAPI.QrWebsocket)
		group.GET("/message_ws", messageWebsocketAPI.MessageWSController)
	}
}

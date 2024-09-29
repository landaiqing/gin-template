package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var clientApi = controller.Controller.ClientController

func ClientRouter(router *gin.RouterGroup) {
	router.GET("/client/generate_client_id", clientApi.GenerateClientId)
}

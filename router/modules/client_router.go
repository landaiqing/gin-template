package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var clientApi = api.Api.ClientApi

func ClientRouter(router *gin.RouterGroup) {
	router.GET("/client/generate_client_id", clientApi.GenerateClientId)
}

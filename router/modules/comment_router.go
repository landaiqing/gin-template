package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api"
)

var commonApi = api.Api.CommonApi

func CommentRouter(router *gin.RouterGroup) {
	router.POST("/auth/comment/submit", commonApi.CommentSubmit)
}

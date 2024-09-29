package modules

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/controller"
)

var commonApi = controller.Controller.CommonController

func CommentRouter(router *gin.RouterGroup) {
	router.POST("/auth/comment/submit", commonApi.CommentSubmit)
	router.POST("/auth/reply/submit", commonApi.ReplySubmit)
	router.POST("/auth/comment/list", commonApi.CommentList)
	router.POST("/auth/reply/list", commonApi.ReplyList)
	router.POST("/auth/reply/reply/submit", commonApi.ReplyReplySubmit)
	router.POST("/auth/comment/like", commonApi.CommentLikes)
	router.POST("/auth/comment/cancel_like", commonApi.CancelCommentLikes)
}

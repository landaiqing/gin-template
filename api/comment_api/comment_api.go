package comment_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api/comment_api/dto"
	"schisandra-cloud-album/common/result"
)

// CommentSubmit 提交评论
func (CommentAPI) CommentSubmit(c *gin.Context) {
	commentRequest := dto.CommentRequest{}
	err := c.ShouldBindJSON(&commentRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	content := commentRequest.Content
	images := commentRequest.Images
	if content == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "RequestParamsNotEmpty"), c)
		return
	}
	if len(images) > 5 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "TooManyImages"), c)
		return
	}
}

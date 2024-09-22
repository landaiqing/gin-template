package comment_api

import (
	"context"
	"errors"
	"github.com/acmestack/gorm-plus/gplus"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schisandra-cloud-album/api/comment_api/dto"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"time"
)

// CommentSubmit 提交评论
// @Summary 提交评论
// @Description 提交评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_request body dto.CommentRequest true "评论请求"
// @Router /auth/comment/submit [post]
func (CommentAPI) CommentSubmit(c *gin.Context) {
	commentRequest := dto.CommentRequest{}
	if err := c.ShouldBindJSON(&commentRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	if commentRequest.Content == "" || commentRequest.UserID == "" || commentRequest.TopicId == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	if len(commentRequest.Images) > 3 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "TooManyImages"), c)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		global.LOG.Errorln("user-agent is empty")
		return
	}
	ua := useragent.New(userAgent)

	ip := utils.GetClientIP(c)
	location, err := global.IP2Location.SearchByStr(ip)
	if err != nil {
		global.LOG.Errorln(err)
		return
	}
	location = utils.RemoveZeroAndAdjust(location)

	browser, _ := ua.Browser()
	operatingSystem := ua.OS()
	isAuthor := 0
	if commentRequest.UserID == commentRequest.Author {
		isAuthor = 1
	}
	commentReply := model.ScaCommentReply{
		Content:         commentRequest.Content,
		UserId:          commentRequest.UserID,
		TopicId:         commentRequest.TopicId,
		TopicType:       enum.CommentTopicType,
		CommentType:     enum.COMMENT,
		Author:          isAuthor,
		CommentIp:       ip,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
	}

	if err = commentReplyService.CreateCommentReply(&commentReply); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}

	if len(commentRequest.Images) > 0 {
		commentImages := CommentImages{
			TopicId:   commentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    commentRequest.UserID,
			Images:    commentRequest.Images,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
			global.LOG.Errorln(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			return
		}
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitSuccess"), c)
	return
}

// ReplySubmit 提交回复
// @Summary 提交回复
// @Description 提交回复
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_comment_request body dto.ReplyCommentRequest true "回复评论请求"
// @Router /auth/reply/submit [post]
func (CommentAPI) ReplySubmit(c *gin.Context) {
	replyCommentRequest := dto.ReplyCommentRequest{}
	if err := c.ShouldBindJSON(&replyCommentRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	if replyCommentRequest.Content == "" ||
		replyCommentRequest.UserID == "" ||
		replyCommentRequest.TopicId == "" ||
		replyCommentRequest.ReplyId == "" ||
		replyCommentRequest.ReplyUser == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	if len(replyCommentRequest.Images) > 3 {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "TooManyImages"), c)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		global.LOG.Errorln("user-agent is empty")
		return
	}
	ua := useragent.New(userAgent)

	ip := utils.GetClientIP(c)
	location, err := global.IP2Location.SearchByStr(ip)
	if err != nil {
		global.LOG.Errorln(err)
		return
	}
	location = utils.RemoveZeroAndAdjust(location)

	browser, _ := ua.Browser()
	operatingSystem := ua.OS()
	isAuthor := 0
	if replyCommentRequest.UserID == replyCommentRequest.Author {
		isAuthor = 1
	}
	commentReply := model.ScaCommentReply{
		Content:         replyCommentRequest.Content,
		UserId:          replyCommentRequest.UserID,
		TopicId:         replyCommentRequest.TopicId,
		TopicType:       enum.CommentTopicType,
		CommentType:     enum.REPLY,
		ReplyId:         replyCommentRequest.ReplyId,
		ReplyUser:       replyCommentRequest.ReplyUser,
		Author:          isAuthor,
		CommentIp:       ip,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
	}

	if err = commentReplyService.CreateCommentReply(&commentReply); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}

	if len(replyCommentRequest.Images) > 0 {
		commentImages := CommentImages{
			TopicId:   replyCommentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    replyCommentRequest.UserID,
			Images:    replyCommentRequest.Images,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
			global.LOG.Errorln(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			return
		}
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitSuccess"), c)
	return
}

// CommentList 获取评论列表
// @Summary 获取评论列表
// @Description 获取评论列表
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_list_request body dto.CommentListRequest true "评论列表请求"
// @Router /auth/comment/list [post]
func (CommentAPI) CommentList(c *gin.Context) {
	commentListRequest := dto.CommentListRequest{}
	err := c.ShouldBindJSON(&commentListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	if commentListRequest.TopicId == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](commentListRequest.Page, commentListRequest.Size)
	query.Eq(&u.TopicId, commentListRequest.TopicId).OrderByDesc(&u.Likes)
	page, _ = gplus.SelectPage(page, query)

	var commentsWithImages []CommentData

	for _, user := range page.Records {
		// 获取评论图片
		commentImages := CommentImages{}
		wrong := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").FindOne(context.Background(), bson.M{"comment_id": user.Id}).Decode(&commentImages)
		if wrong != nil && !errors.Is(wrong, mongo.ErrNoDocuments) {
			global.LOG.Errorln(wrong)
		}
		commentsWithImages = append(commentsWithImages, CommentData{
			Comment: *user,
			Images:  commentImages.Images,
		})
	}
	response := CommentResponse{
		Comments: commentsWithImages,
		Size:     page.Size,
		Current:  page.Current,
		Total:    page.Total,
	}
	result.OkWithData(response, c)
	return
}

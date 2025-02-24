package comment_controller

import (
	"time"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service/impl"
	"schisandra-cloud-album/utils"
)

type CommentController struct{}

var commentReplyService = impl.CommentReplyServiceImpl{}

// CommentSubmit 提交评论
// @Summary 提交评论
// @Description 提交评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_request body CommentRequest true "评论请求"
// @Router /auth/comment/submit [post]
func (CommentController) CommentSubmit(c *gin.Context) {
	commentRequest := CommentRequest{}
	if err := c.ShouldBindJSON(&commentRequest); err != nil {
		return
	}
	// 验证校验
	res := utils.CheckSlideData(commentRequest.Point, commentRequest.Key)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
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
	uid := utils.GetSession(c, constant.SessionKey).UID
	if uid == commentRequest.Author {
		isAuthor = 1
	}
	xssFilterContent := utils.XssFilter(commentRequest.Content)
	if xssFilterContent == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	commentContent := global.SensitiveManager.Replace(xssFilterContent, '*')

	commentReply := model.ScaCommentReply{
		Content:         commentContent,
		UserId:          uid,
		TopicId:         commentRequest.TopicId,
		TopicType:       enum.CommentTopicType,
		CommentType:     enum.COMMENT,
		Author:          isAuthor,
		CommentIp:       ip,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
		Agent:           userAgent,
	}
	commentId, response := commentReplyService.SubmitCommentService(&commentReply, commentRequest.TopicId, uid, commentRequest.Images)
	if !response {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	responseData := model.ScaCommentReply{
		Id:              commentId,
		Content:         commentContent,
		UserId:          uid,
		TopicId:         commentRequest.TopicId,
		Author:          isAuthor,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
		CreatedTime:     time.Now(),
	}
	result.OkWithData(responseData, c)
	return
}

// ReplySubmit 提交回复
// @Summary 提交回复
// @Description 提交回复
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_comment_request body ReplyCommentRequest true "回复评论请求"
// @Router /auth/reply/submit [post]
func (CommentController) ReplySubmit(c *gin.Context) {
	replyCommentRequest := ReplyCommentRequest{}
	if err := c.ShouldBindJSON(&replyCommentRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	// 验证校验
	res := utils.CheckSlideData(replyCommentRequest.Point, replyCommentRequest.Key)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
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
	uid := utils.GetSession(c, constant.SessionKey).UID
	if uid == replyCommentRequest.Author {
		isAuthor = 1
	}
	xssFilterContent := utils.XssFilter(replyCommentRequest.Content)
	if xssFilterContent == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	commentContent := global.SensitiveManager.Replace(xssFilterContent, '*')
	commentReply := model.ScaCommentReply{
		Content:         commentContent,
		UserId:          uid,
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
		Agent:           userAgent,
	}
	commentReplyId, response := commentReplyService.SubmitCommentService(&commentReply, replyCommentRequest.TopicId, uid, replyCommentRequest.Images)
	if !response {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	responseData := model.ScaCommentReply{
		Id:              commentReplyId,
		Content:         commentContent,
		UserId:          uid,
		TopicId:         replyCommentRequest.TopicId,
		ReplyId:         replyCommentRequest.ReplyId,
		ReplyUser:       replyCommentRequest.ReplyUser,
		Author:          isAuthor,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
		CreatedTime:     time.Now(),
	}
	result.OkWithData(responseData, c)
	return
}

// ReplyReplySubmit 提交回复的回复
// @Summary 提交回复的回复
// @Description 提交回复的回复
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_reply_request body ReplyReplyRequest true "回复回复请求"
// @Router /auth/reply/reply/submit [post]
func (CommentController) ReplyReplySubmit(c *gin.Context) {
	replyReplyRequest := ReplyReplyRequest{}
	if err := c.ShouldBindJSON(&replyReplyRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	// 验证校验
	res := utils.CheckSlideData(replyReplyRequest.Point, replyReplyRequest.Key)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
		return
	}
	if len(replyReplyRequest.Images) > 3 {
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
	uid := utils.GetSession(c, constant.SessionKey).UID
	if uid == replyReplyRequest.Author {
		isAuthor = 1
	}
	xssFilterContent := utils.XssFilter(replyReplyRequest.Content)
	if xssFilterContent == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	commentContent := global.SensitiveManager.Replace(xssFilterContent, '*')
	commentReply := model.ScaCommentReply{
		Content:         commentContent,
		UserId:          uid,
		TopicId:         replyReplyRequest.TopicId,
		TopicType:       enum.CommentTopicType,
		CommentType:     enum.REPLY,
		ReplyTo:         replyReplyRequest.ReplyTo,
		ReplyId:         replyReplyRequest.ReplyId,
		ReplyUser:       replyReplyRequest.ReplyUser,
		Author:          isAuthor,
		CommentIp:       ip,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
		Agent:           userAgent,
	}
	commentReplyReplyId, response := commentReplyService.SubmitCommentService(&commentReply, replyReplyRequest.TopicId, uid, replyReplyRequest.Images)
	if !response {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		return
	}
	responseData := model.ScaCommentReply{
		Id:              commentReplyReplyId,
		Content:         commentContent,
		UserId:          uid,
		TopicId:         replyReplyRequest.TopicId,
		ReplyTo:         replyReplyRequest.ReplyTo,
		ReplyId:         replyReplyRequest.ReplyId,
		ReplyUser:       replyReplyRequest.ReplyUser,
		Author:          isAuthor,
		Location:        location,
		Browser:         browser,
		OperatingSystem: operatingSystem,
		CreatedTime:     time.Now(),
	}
	result.OkWithData(responseData, c)
	return
}

// CommentList 获取评论列表
// @Summary 获取评论列表
// @Description 获取评论列表
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_list_request body CommentListRequest true "评论列表请求"
// @Router /auth/comment/list [post]
func (CommentController) CommentList(c *gin.Context) {
	commentListRequest := CommentListRequest{}
	err := c.ShouldBindJSON(&commentListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	uid := utils.GetSession(c, constant.SessionKey).UID
	response := commentReplyService.GetCommentListService(uid, commentListRequest.TopicId, commentListRequest.Page, commentListRequest.Size, commentListRequest.IsHot)
	result.OkWithData(response, c)
	return
}

// ReplyList 获取回复列表
// @Summary 获取回复列表
// @Description 获取回复列表
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_list_request body ReplyListRequest true "回复列表请求"
// @Router /auth/reply/list [post]
func (CommentController) ReplyList(c *gin.Context) {
	replyListRequest := ReplyListRequest{}
	err := c.ShouldBindJSON(&replyListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	uid := utils.GetSession(c, constant.SessionKey).UID
	response := commentReplyService.GetCommentReplyListService(uid, replyListRequest.TopicId, replyListRequest.CommentId, replyListRequest.Page, replyListRequest.Size)
	result.OkWithData(response, c)
	return
}

// CommentLikes 点赞评论
// @Summary 点赞评论
// @Description 点赞评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_like_request body CommentLikeRequest true "点赞请求"
// @Router /auth/comment/like [post]
func (CommentController) CommentLikes(c *gin.Context) {
	likeRequest := CommentLikeRequest{}
	err := c.ShouldBindJSON(&likeRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	uid := utils.GetSession(c, constant.SessionKey).UID
	res := commentReplyService.CommentLikeService(uid, likeRequest.CommentId, likeRequest.TopicId)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentLikeFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentLikeSuccess"), c)
	return
}

// CancelCommentLikes 取消点赞评论
// @Summary 取消点赞评论
// @Description 取消点赞评论
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param comment_like_request body CommentLikeRequest true "取消点赞请求"
// @Router /auth/comment/cancel_like [post]
func (CommentController) CancelCommentLikes(c *gin.Context) {
	cancelLikeRequest := CommentLikeRequest{}
	if err := c.ShouldBindJSON(&cancelLikeRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	uid := utils.GetSession(c, constant.SessionKey).UID
	res := commentReplyService.CommentDislikeService(uid, cancelLikeRequest.CommentId, cancelLikeRequest.TopicId)
	if !res {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentDislikeFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentDislikeSuccess"), c)
	return
}

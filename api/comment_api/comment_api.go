package comment_api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"regexp"
	"schisandra-cloud-album/api/comment_api/dto"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"strconv"
	"strings"
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
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
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
		Agent:           userAgent,
	}

	if err = commentReplyService.CreateCommentReply(&commentReply); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	if len(commentRequest.Images) > 0 {

		var imagesData [][]byte
		for _, img := range commentRequest.Images {
			re := regexp.MustCompile(`^data:image/\w+;base64,`)
			imgWithoutPrefix := re.ReplaceAllString(img, "")
			data, err := base64ToBytes(imgWithoutPrefix)
			if err != nil {
				global.LOG.Errorln(err)
				tx.Rollback()
				return
			}
			imagesData = append(imagesData, data)
		}

		commentImages := CommentImages{
			TopicId:   commentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    commentRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
			global.LOG.Errorln(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitSuccess"), c)
	return
}

// base64ToBytes 将base64字符串转换为字节数组
func base64ToBytes(base64Str string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str))
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("failed to decode base64 string")
	}
	return data, nil
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
		strconv.FormatInt(replyCommentRequest.ReplyId, 10) == "" ||
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
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
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
		Agent:           userAgent,
	}

	if err = commentReplyService.CreateCommentReply(&commentReply); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}
	err = commentReplyService.UpdateCommentReplyCount(replyCommentRequest.ReplyId)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	if len(replyCommentRequest.Images) > 0 {

		var imagesData [][]byte
		for _, img := range replyCommentRequest.Images {
			re := regexp.MustCompile(`^data:image/\w+;base64,`)
			imgWithoutPrefix := re.ReplaceAllString(img, "")
			data, err := base64ToBytes(imgWithoutPrefix)
			if err != nil {
				global.LOG.Errorln(err)
				tx.Rollback()
				return
			}
			imagesData = append(imagesData, data)
		}
		commentImages := CommentImages{
			TopicId:   replyCommentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    replyCommentRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
			global.LOG.Errorln(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitSuccess"), c)
	return
}

// ReplyReplySubmit 提交回复的回复
// @Summary 提交回复的回复
// @Description 提交回复的回复
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_reply_request body dto.ReplyReplyRequest true "回复回复请求"
// @Router /auth/reply/reply/submit [post]
func (CommentAPI) ReplyReplySubmit(c *gin.Context) {
	replyReplyRequest := dto.ReplyReplyRequest{}
	if err := c.ShouldBindJSON(&replyReplyRequest); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	if replyReplyRequest.Content == "" ||
		replyReplyRequest.UserID == "" ||
		replyReplyRequest.TopicId == "" ||
		replyReplyRequest.ReplyTo == 0 ||
		replyReplyRequest.ReplyId == 0 ||
		replyReplyRequest.ReplyUser == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
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
	if replyReplyRequest.UserID == replyReplyRequest.Author {
		isAuthor = 1
	}
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	commentReply := model.ScaCommentReply{
		Content:         replyReplyRequest.Content,
		UserId:          replyReplyRequest.UserID,
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

	if err = commentReplyService.CreateCommentReply(&commentReply); err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}
	err = commentReplyService.UpdateCommentReplyCount(replyReplyRequest.ReplyId)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	if len(replyReplyRequest.Images) > 0 {

		var imagesData [][]byte
		for _, img := range replyReplyRequest.Images {
			re := regexp.MustCompile(`^data:image/\w+;base64,`)
			imgWithoutPrefix := re.ReplaceAllString(img, "")
			data, err := base64ToBytes(imgWithoutPrefix)
			if err != nil {
				global.LOG.Errorln(err)
				tx.Rollback()
				return
			}
			imagesData = append(imagesData, data)
		}
		commentImages := CommentImages{
			TopicId:   replyReplyRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    replyReplyRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
			global.LOG.Errorln(err)
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
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
	query.Eq(&u.TopicId, commentListRequest.TopicId).Eq(&u.CommentType, enum.COMMENT).OrderByDesc(&u.Likes)
	page, _ = gplus.SelectPage(page, query)

	var commentsWithImages []CommentContent

	for _, comment := range page.Records {
		// 获取评论图片
		commentImages := CommentImages{}
		wrong := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").FindOne(context.Background(), bson.M{"comment_id": comment.Id}).Decode(&commentImages)
		if wrong != nil && !errors.Is(wrong, mongo.ErrNoDocuments) {
			global.LOG.Errorln(wrong)
		}
		// 将图片转换为base64
		var imagesBase64 []string
		for _, img := range commentImages.Images {
			mimeType := getMimeType(img) // 动态获取 MIME 类型
			base64Img := base64.StdEncoding.EncodeToString(img)
			base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
			imagesBase64 = append(imagesBase64, base64WithPrefix)
		}
		// 组装评论数据
		queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
		queryUser.Eq(&n.UID, comment.UserId)
		userInfo, _ := gplus.SelectOne(queryUser)
		commentsWithImages = append(commentsWithImages,
			CommentContent{
				Avatar:          *userInfo.Avatar,
				NickName:        *userInfo.Nickname,
				Id:              comment.Id,
				UserId:          comment.UserId,
				TopicId:         comment.TopicId,
				Dislikes:        comment.Dislikes,
				Content:         comment.Content,
				ReplyCount:      comment.ReplyCount,
				Likes:           comment.Likes,
				CreatedTime:     comment.CreatedTime,
				Author:          comment.Author,
				Location:        comment.Location,
				Browser:         comment.Browser,
				OperatingSystem: comment.OperatingSystem,
				Images:          imagesBase64,
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

func getMimeType(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream" // 默认类型
	}

	// 判断 JPEG
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// 判断 PNG
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// 判断 GIF
	if len(data) >= 6 && data[0] == 'G' && data[1] == 'I' && data[2] == 'F' {
		return "image/gif"
	}

	return "application/octet-stream" // 默认类型
}

// ReplyList 获取回复列表
// @Summary 获取回复列表
// @Description 获取回复列表
// @Tags 评论
// @Accept  json
// @Produce  json
// @Param reply_list_request body dto.ReplyListRequest true "回复列表请求"
// @Router /auth/reply/list [post]
func (CommentAPI) ReplyList(c *gin.Context) {
	replyListRequest := dto.ReplyListRequest{}
	err := c.ShouldBindJSON(&replyListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	if replyListRequest.TopicId == "" || strconv.FormatInt(replyListRequest.CommentId, 10) == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](replyListRequest.Page, replyListRequest.Size)
	query.Eq(&u.TopicId, replyListRequest.TopicId).Eq(&u.ReplyId, replyListRequest.CommentId).Eq(&u.CommentType, enum.REPLY).OrderByDesc(&u.Likes)
	page, _ = gplus.SelectPage(page, query)

	var commentsWithImages []CommentContent

	for _, comment := range page.Records {
		// 获取评论图片
		commentImages := CommentImages{}
		wrong := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").FindOne(context.Background(), bson.M{"comment_id": comment.Id}).Decode(&commentImages)
		if wrong != nil && !errors.Is(wrong, mongo.ErrNoDocuments) {
			global.LOG.Errorln(wrong)
		}
		// 将图片转换为base64
		var imagesBase64 []string
		for _, img := range commentImages.Images {
			mimeType := getMimeType(img) // 动态获取 MIME 类型
			base64Img := base64.StdEncoding.EncodeToString(img)
			base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
			imagesBase64 = append(imagesBase64, base64WithPrefix)
		}
		// 查询评论用户信息
		queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
		queryUser.Eq(&n.UID, comment.UserId)
		userInfo, _ := gplus.SelectOne(queryUser)
		// 查询回复用户信息
		queryReplyUser, m := gplus.NewQuery[model.ScaAuthUser]()
		queryReplyUser.Eq(&m.UID, comment.ReplyUser)
		replyUserInfo, _ := gplus.SelectOne(queryReplyUser)
		commentsWithImages = append(commentsWithImages,
			CommentContent{
				Avatar:          *userInfo.Avatar,
				NickName:        *userInfo.Nickname,
				Id:              comment.Id,
				UserId:          comment.UserId,
				TopicId:         comment.TopicId,
				Dislikes:        comment.Dislikes,
				Content:         comment.Content,
				ReplyUsername:   *replyUserInfo.Nickname,
				ReplyCount:      comment.ReplyCount,
				Likes:           comment.Likes,
				CreatedTime:     comment.CreatedTime,
				ReplyUser:       comment.ReplyUser,
				ReplyId:         comment.ReplyId,
				ReplyTo:         comment.ReplyTo,
				Author:          comment.Author,
				Location:        comment.Location,
				Browser:         comment.Browser,
				OperatingSystem: comment.OperatingSystem,
				Images:          imagesBase64,
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

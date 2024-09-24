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
	// 使用 goroutine 进行异步评论保存
	errCh := make(chan error, 2)
	go func() {
		errCh <- commentReplyService.CreateCommentReply(&commentReply)
	}()

	// 等待评论回复的创建
	if err = <-errCh; err != nil {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	// 处理图片异步上传
	if len(commentRequest.Images) > 0 {
		imagesDataCh := make(chan [][]byte)
		go func() {
			imagesData, err := processImages(commentRequest.Images)
			if err != nil {
				global.LOG.Errorln(err)
				imagesDataCh <- nil // 发送失败信号
				return
			}
			imagesDataCh <- imagesData // 发送处理成功的数据
		}()

		imagesData := <-imagesDataCh
		if imagesData == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}

		commentImages := CommentImages{
			TopicId:   commentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    commentRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 使用 goroutine 进行异步图片保存
		go func() {
			if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
				global.LOG.Errorln(err)
			}
		}()
	}
	tx.Commit()
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
	// 使用 goroutine 进行异步评论保存
	errCh := make(chan error)
	go func() {

		errCh <- commentReplyService.CreateCommentReply(&commentReply)
	}()
	go func() {

		errCh <- commentReplyService.UpdateCommentReplyCount(replyCommentRequest.ReplyId)
	}()
	// 等待评论回复的创建
	if err = <-errCh; err != nil {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	// 处理图片异步上传
	if len(replyCommentRequest.Images) > 0 {
		imagesDataCh := make(chan [][]byte)
		go func() {
			imagesData, err := processImages(replyCommentRequest.Images)
			if err != nil {
				global.LOG.Errorln(err)
				imagesDataCh <- nil // 发送失败信号
				return
			}
			imagesDataCh <- imagesData // 发送处理成功的数据
		}()

		imagesData := <-imagesDataCh
		if imagesData == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}

		commentImages := CommentImages{
			TopicId:   replyCommentRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    replyCommentRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 使用 goroutine 进行异步图片保存
		go func() {
			if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
				global.LOG.Errorln(err)
			}
		}()
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

	errCh := make(chan error, 2)
	go func() {
		errCh <- commentReplyService.CreateCommentReply(&commentReply)
	}()
	go func() {
		errCh <- commentReplyService.UpdateCommentReplyCount(replyReplyRequest.ReplyId)
	}()

	if err = <-errCh; err != nil {
		global.LOG.Errorln(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
		tx.Rollback()
		return
	}

	if len(replyReplyRequest.Images) > 0 {
		imagesDataCh := make(chan [][]byte)
		go func() {
			imagesData, err := processImages(replyReplyRequest.Images)
			if err != nil {
				global.LOG.Errorln(err)
				imagesDataCh <- nil
				return
			}
			imagesDataCh <- imagesData
		}()

		imagesData := <-imagesDataCh
		if imagesData == nil {
			result.FailWithMessage(ginI18n.MustGetMessage(c, "CommentSubmitFailed"), c)
			tx.Rollback()
			return
		}

		commentImages := CommentImages{
			TopicId:   replyReplyRequest.TopicId,
			CommentId: commentReply.Id,
			UserId:    replyReplyRequest.UserID,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 处理图片保存
		go func() {
			if _, err = global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
				global.LOG.Errorln(err)
			}
		}()
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
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](commentListRequest.Page, commentListRequest.Size)
	query.Eq(&u.TopicId, commentListRequest.TopicId).Eq(&u.CommentType, enum.COMMENT).OrderByDesc(&u.Likes)
	page, _ = gplus.SelectPage(page, query)

	userIds := make([]string, 0, len(page.Records))
	for _, comment := range page.Records {
		userIds = append(userIds, comment.UserId)
	}

	queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
	queryUser.In(&n.UID, userIds)
	userInfos, _ := gplus.SelectList(queryUser)

	userInfoMap := make(map[string]model.ScaAuthUser, len(userInfos))
	for _, userInfo := range userInfos {
		userInfoMap[*userInfo.UID] = *userInfo
	}

	commentChannel := make(chan CommentContent, len(page.Records))
	imagesBase64S := make([][]string, len(page.Records)) // 存储每条评论的图片

	for index, comment := range page.Records {
		wg.Add(1)
		go func(comment model.ScaCommentReply, index int) {
			defer wg.Done()

			// 使用 context 设置超时时间
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 设置超时，2秒
			defer cancel()

			// 获取评论图片并处理
			commentImages := CommentImages{}
			wrong := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").FindOne(ctx, bson.M{"comment_id": comment.Id}).Decode(&commentImages)
			if wrong != nil && !errors.Is(wrong, mongo.ErrNoDocuments) {
				global.LOG.Errorf("Failed to get images for comment ID %s: %v", comment.Id, wrong)
				return
			}

			// 将图片转换为base64
			var imagesBase64 []string
			for _, img := range commentImages.Images {
				mimeType := getMimeType(img)
				base64Img := base64.StdEncoding.EncodeToString(img)
				base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
				imagesBase64 = append(imagesBase64, base64WithPrefix)
			}
			imagesBase64S[index] = imagesBase64 // 保存到切片中

			userInfo := userInfoMap[comment.UserId]
			commentContent := CommentContent{
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
			}
			commentChannel <- commentContent
		}(*comment, index)
	}

	go func() {
		wg.Wait()
		close(commentChannel)
	}()

	var commentsWithImages []CommentContent
	for commentContent := range commentChannel {
		commentsWithImages = append(commentsWithImages, commentContent)
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
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](replyListRequest.Page, replyListRequest.Size)
	query.Eq(&u.TopicId, replyListRequest.TopicId).Eq(&u.ReplyId, replyListRequest.CommentId).Eq(&u.CommentType, enum.REPLY).OrderByDesc(&u.Likes)
	page, _ = gplus.SelectPage(page, query)

	userIds := make([]string, 0, len(page.Records))
	replyUserIds := make([]string, 0, len(page.Records))

	// 收集用户 ID 和回复用户 ID
	for _, comment := range page.Records {
		userIds = append(userIds, comment.UserId)
		if comment.ReplyUser != "" {
			replyUserIds = append(replyUserIds, comment.ReplyUser)
		}
	}

	// 查询评论用户信息
	queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
	queryUser.In(&n.UID, userIds)
	userInfos, _ := gplus.SelectList(queryUser)

	userInfoMap := make(map[string]model.ScaAuthUser, len(userInfos))
	for _, userInfo := range userInfos {
		userInfoMap[*userInfo.UID] = *userInfo
	}

	// 查询回复用户信息
	replyUserInfoMap := make(map[string]model.ScaAuthUser)
	if len(replyUserIds) > 0 {
		queryReplyUser, m := gplus.NewQuery[model.ScaAuthUser]()
		queryReplyUser.In(&m.UID, replyUserIds)
		replyUserInfos, _ := gplus.SelectList(queryReplyUser)

		for _, replyUserInfo := range replyUserInfos {
			replyUserInfoMap[*replyUserInfo.UID] = *replyUserInfo
		}
	}

	replyChannel := make(chan CommentContent, len(page.Records)) // 使用通道传递回复内容
	imagesBase64S := make([][]string, len(page.Records))         // 存储每条回复的图片

	for index, reply := range page.Records {
		wg.Add(1)
		go func(reply model.ScaCommentReply, index int) {
			defer wg.Done()

			// 使用 context 设置超时时间
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 设置超时，2秒
			defer cancel()

			// 获取回复图片并处理
			replyImages := CommentImages{}
			wrong := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").FindOne(ctx, bson.M{"comment_id": reply.Id}).Decode(&replyImages)
			if wrong != nil && !errors.Is(wrong, mongo.ErrNoDocuments) {
				global.LOG.Errorf("Failed to get images for reply ID %s: %v", reply.Id, wrong)
				return
			}

			// 将图片转换为base64
			var imagesBase64 []string
			for _, img := range replyImages.Images {
				mimeType := getMimeType(img)
				base64Img := base64.StdEncoding.EncodeToString(img)
				base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
				imagesBase64 = append(imagesBase64, base64WithPrefix)
			}
			imagesBase64S[index] = imagesBase64 // 保存到切片中

			userInfo := userInfoMap[reply.UserId]
			replyUserInfo := replyUserInfoMap[reply.ReplyUser]
			commentContent := CommentContent{
				Avatar:          *userInfo.Avatar,
				NickName:        *userInfo.Nickname,
				Id:              reply.Id,
				UserId:          reply.UserId,
				TopicId:         reply.TopicId,
				Dislikes:        reply.Dislikes,
				Content:         reply.Content,
				ReplyUsername:   *replyUserInfo.Nickname,
				ReplyCount:      reply.ReplyCount,
				Likes:           reply.Likes,
				CreatedTime:     reply.CreatedTime,
				ReplyUser:       reply.ReplyUser,
				ReplyId:         reply.ReplyId,
				ReplyTo:         reply.ReplyTo,
				Author:          reply.Author,
				Location:        reply.Location,
				Browser:         reply.Browser,
				OperatingSystem: reply.OperatingSystem,
				Images:          imagesBase64,
			}
			replyChannel <- commentContent // 发送到通道
		}(*reply, index)
	}

	go func() {
		wg.Wait()
		close(replyChannel) // 关闭通道
	}()

	var repliesWithImages []CommentContent
	for replyContent := range replyChannel {
		repliesWithImages = append(repliesWithImages, replyContent) // 从通道获取回复内容
	}

	response := CommentResponse{
		Comments: repliesWithImages,
		Size:     page.Size,
		Current:  page.Current,
		Total:    page.Total,
	}
	result.OkWithData(response, c)
	return
}

// CommentLikes 点赞评论
func (CommentAPI) CommentLikes(c *gin.Context) {
	likeRequest := dto.CommentLikeRequest{}
	err := c.ShouldBindJSON(&likeRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
}

func (CommentAPI) CommentDislikes(c *gin.Context) {

}

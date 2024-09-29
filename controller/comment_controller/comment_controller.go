package comment_controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/mq"
	"schisandra-cloud-album/utils"
	"time"
)

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
		errCh <- commentReplyService.CreateCommentReplyService(&commentReply)
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

		errCh <- commentReplyService.CreateCommentReplyService(&commentReply)
	}()
	go func() {

		errCh <- commentReplyService.UpdateCommentReplyCountService(replyCommentRequest.ReplyId)
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
		errCh <- commentReplyService.CreateCommentReplyService(&commentReply)
	}()
	go func() {
		errCh <- commentReplyService.UpdateCommentReplyCountService(replyReplyRequest.ReplyId)
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
// @Param comment_list_request body CommentListRequest true "评论列表请求"
// @Router /auth/comment/list [post]
func (CommentController) CommentList(c *gin.Context) {
	commentListRequest := CommentListRequest{}
	err := c.ShouldBindJSON(&commentListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	// 查询评论列表
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](commentListRequest.Page, commentListRequest.Size)
	if commentListRequest.IsHot {
		query.OrderByDesc(&u.CommentOrder).OrderByDesc(&u.Likes).OrderByDesc(&u.ReplyCount)
	} else {
		query.OrderByDesc(&u.CommentOrder).OrderByDesc(&u.CreatedTime)
	}
	query.Eq(&u.TopicId, commentListRequest.TopicId).Eq(&u.CommentType, enum.COMMENT)
	page, pageDB := gplus.SelectPage(page, query)
	if pageDB.Error != nil {
		global.LOG.Errorln(pageDB.Error)
		return
	}
	if len(page.Records) == 0 {
		result.OkWithData(CommentResponse{Comments: []CommentContent{}, Size: page.Size, Current: page.Current, Total: page.Total}, c)
		return
	}

	userIds := make([]string, 0, len(page.Records))
	commentIds := make([]int64, 0, len(page.Records))
	for _, comment := range page.Records {
		userIds = append(userIds, comment.UserId)
		commentIds = append(commentIds, comment.Id)
	}

	// 结果存储
	userInfoMap := make(map[string]model.ScaAuthUser)
	likeMap := make(map[int64]bool)
	commentImagesMap := make(map[int64]CommentImages)

	// 使用 WaitGroup 等待协程完成
	wg.Add(3)

	// 查询评论用户信息
	go func() {
		defer wg.Done()
		queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
		queryUser.Select(&n.UID, &n.Avatar, &n.Nickname).In(&n.UID, userIds)
		userInfos, userInfosDB := gplus.SelectList(queryUser)
		if userInfosDB.Error != nil {
			global.LOG.Errorln(userInfosDB.Error)
			return
		}
		for _, userInfo := range userInfos {
			userInfoMap[*userInfo.UID] = *userInfo
		}
	}()

	// 查询评论点赞状态
	go func() {
		defer wg.Done()
		if len(page.Records) > 0 {
			queryLike, l := gplus.NewQuery[model.ScaCommentLikes]()
			queryLike.Eq(&l.TopicId, commentListRequest.TopicId).Eq(&l.UserId, commentListRequest.UserID).In(&l.CommentId, commentIds)
			likes, likesDB := gplus.SelectList(queryLike)
			if likesDB.Error != nil {
				global.LOG.Errorln(likesDB.Error)
				return
			}
			for _, like := range likes {
				likeMap[like.CommentId] = true
			}
		}
	}()

	// 查询评论图片信息
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置超时，2秒
		defer cancel()

		cursor, err := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").Find(ctx, bson.M{"comment_id": bson.M{"$in": commentIds}})
		if err != nil {
			global.LOG.Errorf("Failed to get images for comments: %v", err)
			return
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
			if err != nil {
				return
			}
		}(cursor, ctx)

		for cursor.Next(ctx) {
			var commentImages CommentImages
			if err = cursor.Decode(&commentImages); err != nil {
				global.LOG.Errorf("Failed to decode comment images: %v", err)
				continue
			}
			commentImagesMap[commentImages.CommentId] = commentImages
		}
	}()

	// 等待所有查询完成
	wg.Wait()
	commentChannel := make(chan CommentContent, len(page.Records))

	for _, comment := range page.Records {
		wg.Add(1)
		go func(comment model.ScaCommentReply) {
			defer wg.Done()
			// 将图片转换为base64
			var imagesBase64 []string
			if imgData, ok := commentImagesMap[comment.Id]; ok {
				// 将图片转换为base64
				for _, img := range imgData.Images {
					mimeType := getMimeType(img)
					base64Img := base64.StdEncoding.EncodeToString(img)
					base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
					imagesBase64 = append(imagesBase64, base64WithPrefix)
				}
			}

			userInfo, exist := userInfoMap[comment.UserId]
			if !exist {
				global.LOG.Errorf("Failed to get user info for comment: %s", comment.UserId)
				return
			}
			commentContent := CommentContent{
				Avatar:          *userInfo.Avatar,
				NickName:        *userInfo.Nickname,
				Id:              comment.Id,
				UserId:          comment.UserId,
				TopicId:         comment.TopicId,
				Content:         comment.Content,
				ReplyCount:      comment.ReplyCount,
				Likes:           comment.Likes,
				CreatedTime:     comment.CreatedTime,
				Author:          comment.Author,
				Location:        comment.Location,
				Browser:         comment.Browser,
				OperatingSystem: comment.OperatingSystem,
				Images:          imagesBase64,
				IsLiked:         likeMap[comment.Id],
			}
			commentChannel <- commentContent
		}(*comment)
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
// @Param reply_list_request body ReplyListRequest true "回复列表请求"
// @Router /auth/reply/list [post]
func (CommentController) ReplyList(c *gin.Context) {
	replyListRequest := ReplyListRequest{}
	err := c.ShouldBindJSON(&replyListRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](replyListRequest.Page, replyListRequest.Size)
	query.Eq(&u.TopicId, replyListRequest.TopicId).Eq(&u.ReplyId, replyListRequest.CommentId).Eq(&u.CommentType, enum.REPLY).OrderByDesc(&u.Likes).OrderByAsc(&u.CreatedTime)
	page, pageDB := gplus.SelectPage(page, query)
	if pageDB.Error != nil {
		global.LOG.Errorln(pageDB.Error)
		return
	}
	if len(page.Records) == 0 {
		result.OkWithData(CommentResponse{Comments: []CommentContent{}, Size: page.Size, Current: page.Current, Total: page.Total}, c)
		return
	}

	userIdsSet := make(map[string]struct{}) // 使用集合去重用户 ID
	commentIds := make([]int64, 0, len(page.Records))
	// 收集用户 ID 和评论 ID
	for _, comment := range page.Records {
		userIdsSet[comment.UserId] = struct{}{} // 去重
		commentIds = append(commentIds, comment.Id)
		if comment.ReplyUser != "" {
			userIdsSet[comment.ReplyUser] = struct{}{} // 去重
		}
	}
	// 将用户 ID 转换为切片
	userIds := make([]string, 0, len(userIdsSet))
	for userId := range userIdsSet {
		userIds = append(userIds, userId)
	}

	likeMap := make(map[int64]bool, len(page.Records))
	commentImagesMap := make(map[int64]CommentImages)
	userInfoMap := make(map[string]model.ScaAuthUser, len(userIds))

	wg.Add(3)
	go func() {
		defer wg.Done()
		// 查询评论用户信息
		queryUser, n := gplus.NewQuery[model.ScaAuthUser]()
		queryUser.Select(&n.UID, &n.Avatar, &n.Nickname).In(&n.UID, userIds)
		userInfos, userInfosDB := gplus.SelectList(queryUser)
		if userInfosDB.Error != nil {
			global.LOG.Errorln(userInfosDB.Error)
			return
		}
		for _, userInfo := range userInfos {
			userInfoMap[*userInfo.UID] = *userInfo
		}
	}()

	go func() {
		defer wg.Done()
		// 查询评论点赞状态

		if len(page.Records) > 0 {
			queryLike, l := gplus.NewQuery[model.ScaCommentLikes]()
			queryLike.Eq(&l.TopicId, replyListRequest.TopicId).Eq(&l.UserId, replyListRequest.UserID).In(&l.CommentId, commentIds)
			likes, likesDB := gplus.SelectList(queryLike)
			if likesDB.Error != nil {
				global.LOG.Errorln(likesDB.Error)
				return
			}
			for _, like := range likes {
				likeMap[like.CommentId] = true
			}
		}
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 设置超时，2秒
		defer cancel()
		cursor, err := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").Find(ctx, bson.M{"comment_id": bson.M{"$in": commentIds}})
		if err != nil {
			global.LOG.Errorf("Failed to get images for comments: %v", err)
			return
		}
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			warn := cursor.Close(ctx)
			if warn != nil {
				return
			}
		}(cursor, ctx)

		for cursor.Next(ctx) {
			var commentImages CommentImages
			if e := cursor.Decode(&commentImages); e != nil {
				global.LOG.Errorf("Failed to decode comment images: %v", e)
				continue
			}
			commentImagesMap[commentImages.CommentId] = commentImages
		}
	}()
	wg.Wait()

	replyChannel := make(chan CommentContent, len(page.Records)) // 使用通道传递回复内容

	for _, reply := range page.Records {
		wg.Add(1)
		go func(reply model.ScaCommentReply) {
			defer wg.Done()

			var imagesBase64 []string
			if imgData, ok := commentImagesMap[reply.Id]; ok {
				// 将图片转换为base64
				for _, img := range imgData.Images {
					mimeType := getMimeType(img)
					base64Img := base64.StdEncoding.EncodeToString(img)
					base64WithPrefix := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Img)
					imagesBase64 = append(imagesBase64, base64WithPrefix)
				}
			}
			userInfo, exist := userInfoMap[reply.UserId]
			if !exist {
				global.LOG.Errorf("Failed to get user info for comment: %s", reply.UserId)
				return
			}
			replyUserInfo, e := userInfoMap[reply.ReplyUser]
			if !e {
				global.LOG.Errorf("Failed to get reply user info for comment: %s", reply.ReplyUser)
				return
			}
			commentContent := CommentContent{
				Avatar:          *userInfo.Avatar,
				NickName:        *userInfo.Nickname,
				Id:              reply.Id,
				UserId:          reply.UserId,
				TopicId:         reply.TopicId,
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
				IsLiked:         likeMap[reply.Id],
			}
			replyChannel <- commentContent // 发送到通道
		}(*reply)
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

	mx.Lock()
	defer mx.Unlock()

	likes := model.ScaCommentLikes{
		CommentId: likeRequest.CommentId,
		UserId:    likeRequest.UserID,
		TopicId:   likeRequest.TopicId,
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	res := global.DB.Create(&likes) // 假设这是插入数据库的方法
	if res.Error != nil {
		tx.Rollback()
		global.LOG.Errorln(res.Error)
		return
	}

	// 异步更新点赞计数
	go func() {
		if err = commentReplyService.UpdateCommentLikesCountService(likeRequest.CommentId, likeRequest.TopicId); err != nil {
			global.LOG.Errorln(err)
		}
	}()
	marshal, err := json.Marshal(likes)
	if err != nil {
		global.LOG.Errorln(err)
		return
	}
	mq.CommentLikeProducer(marshal)

	tx.Commit()
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
	mx.Lock()
	defer mx.Unlock()

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query, u := gplus.NewQuery[model.ScaCommentLikes]()
	query.Eq(&u.CommentId, cancelLikeRequest.CommentId).
		Eq(&u.UserId, cancelLikeRequest.UserID).
		Eq(&u.TopicId, cancelLikeRequest.TopicId)

	res := gplus.Delete[model.ScaCommentLikes](query)
	if res.Error != nil {
		tx.Rollback()
		return // 返回错误而非打印
	}

	// 异步更新点赞计数
	go func() {
		if err := commentReplyService.DecrementCommentLikesCountService(cancelLikeRequest.CommentId, cancelLikeRequest.TopicId); err != nil {
			global.LOG.Errorln(err)
		}
	}()
	tx.Commit()
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CommentLikeCancelSuccess"), c)
	return
}

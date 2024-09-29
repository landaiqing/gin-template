package impl

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/utils"
	"sync"
	"time"
)

var commentReplyDao = impl.CommentReplyDaoImpl{}

type CommentReplyServiceImpl struct{}

var wg sync.WaitGroup
var mx sync.Mutex

// CommentImages 评论图片
type CommentImages struct {
	TopicId   string   `json:"topic_id" bson:"topic_id" required:"true"`
	CommentId int64    `json:"comment_id" bson:"comment_id" required:"true"`
	UserId    string   `json:"user_id" bson:"user_id" required:"true"`
	Images    [][]byte `json:"images" bson:"images" required:"true"`
	CreatedAt string   `json:"created_at" bson:"created_at" required:"true"`
}

// CommentContent 评论内容
type CommentContent struct {
	NickName        string    `json:"nickname"`
	Avatar          string    `json:"avatar"`
	Level           int       `json:"level,omitempty"`
	Id              int64     `json:"id"`
	UserId          string    `json:"user_id"`
	TopicId         string    `json:"topic_id"`
	Content         string    `json:"content"`
	ReplyTo         int64     `json:"reply_to,omitempty"`
	ReplyId         int64     `json:"reply_id,omitempty"`
	ReplyUser       string    `json:"reply_user,omitempty"`
	ReplyUsername   string    `json:"reply_username,omitempty"`
	Author          int       `json:"author"`
	Likes           int64     `json:"likes"`
	ReplyCount      int64     `json:"reply_count"`
	CreatedTime     time.Time `json:"created_time"`
	Location        string    `json:"location"`
	Browser         string    `json:"browser"`
	OperatingSystem string    `json:"operating_system"`
	IsLiked         bool      `json:"is_liked" default:"false"`
	Images          []string  `json:"images,omitempty"`
}

// CommentResponse 评论返回值
type CommentResponse struct {
	Size     int              `json:"size"`
	Total    int64            `json:"total"`
	Current  int              `json:"current"`
	Comments []CommentContent `json:"comments"`
}

// SubmitCommentService 提交评论
func (CommentReplyServiceImpl) SubmitCommentService(comment *model.ScaCommentReply, topicId string, uid string, images []string) bool {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	errCh := make(chan error, 2)
	go func() {
		errCh <- commentReplyDao.CreateCommentReply(comment)
	}()
	go func() {
		errCh <- commentReplyDao.UpdateCommentReplyCount(comment.ReplyId)
	}()

	if err := <-errCh; err != nil {
		global.LOG.Errorln(err)
		tx.Rollback()
		return false
	}

	if len(images) > 0 {
		imagesDataCh := make(chan [][]byte)
		go func() {
			imagesData, err := utils.ProcessImages(images)
			if err != nil {
				global.LOG.Errorln(err)
				imagesDataCh <- nil
				return
			}
			imagesDataCh <- imagesData
		}()

		imagesData := <-imagesDataCh
		if imagesData == nil {
			tx.Rollback()
			return false
		}

		commentImages := CommentImages{
			TopicId:   topicId,
			CommentId: comment.Id,
			UserId:    uid,
			Images:    imagesData,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 处理图片保存
		go func() {
			if _, err := global.MongoDB.Database(global.CONFIG.MongoDB.DB).Collection("comment_images").InsertOne(context.Background(), commentImages); err != nil {
				global.LOG.Errorln(err)
			}
		}()
	}
	tx.Commit()
	return true
}

// GetCommentReplyService 获取评论回复
func (CommentReplyServiceImpl) GetCommentReplyService(uid string, topicId string, commentId int64, pageNum int, size int) *CommentResponse {
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](pageNum, size)
	query.Eq(&u.TopicId, topicId).Eq(&u.ReplyId, commentId).Eq(&u.CommentType, enum.REPLY).OrderByDesc(&u.Likes).OrderByAsc(&u.CreatedTime)
	page, pageDB := gplus.SelectPage(page, query)
	if pageDB.Error != nil {
		global.LOG.Errorln(pageDB.Error)
		return nil
	}
	if len(page.Records) == 0 {
		return nil

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
			queryLike.Eq(&l.TopicId, topicId).Eq(&l.UserId, uid).In(&l.CommentId, commentIds)
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
					mimeType := utils.GetMimeType(img)
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
	return &response
}

// GetCommentListService 评论列表
func (CommentReplyServiceImpl) GetCommentListService(uid string, topicId string, pageNum int, size int, isHot bool) *CommentResponse {
	// 查询评论列表
	query, u := gplus.NewQuery[model.ScaCommentReply]()
	page := gplus.NewPage[model.ScaCommentReply](pageNum, size)
	if isHot {
		query.OrderByDesc(&u.CommentOrder).OrderByDesc(&u.Likes).OrderByDesc(&u.ReplyCount)
	} else {
		query.OrderByDesc(&u.CommentOrder).OrderByDesc(&u.CreatedTime)
	}
	query.Eq(&u.TopicId, topicId).Eq(&u.CommentType, enum.COMMENT)
	page, pageDB := gplus.SelectPage(page, query)
	if pageDB.Error != nil {
		global.LOG.Errorln(pageDB.Error)
		return nil
	}
	if len(page.Records) == 0 {
		return nil
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
			queryLike.Eq(&l.TopicId, topicId).Eq(&l.UserId, uid).In(&l.CommentId, commentIds)
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
					mimeType := utils.GetMimeType(img)
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
	return &response
}

// CommentLikeService 评论点赞
func (CommentReplyServiceImpl) CommentLikeService(uid string, commentId int64, topicId string) bool {
	mx.Lock()
	defer mx.Unlock()
	likes := model.ScaCommentLikes{
		CommentId: commentId,
		UserId:    uid,
		TopicId:   topicId,
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
		return false
	}

	// 异步更新点赞计数
	go func() {
		if err := commentReplyDao.UpdateCommentLikesCount(commentId, topicId); err != nil {
			global.LOG.Errorln(err)
		}
	}()
	tx.Commit()
	return true
}

// CommentDislikeService 取消评论点赞
func (CommentReplyServiceImpl) CommentDislikeService(uid string, commentId int64, topicId string) bool {
	mx.Lock()
	defer mx.Unlock()

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query, u := gplus.NewQuery[model.ScaCommentLikes]()
	query.Eq(&u.CommentId, commentId).
		Eq(&u.UserId, uid).
		Eq(&u.TopicId, topicId)

	res := gplus.Delete[model.ScaCommentLikes](query)
	if res.Error != nil {
		tx.Rollback()
		return false
	}

	// 异步更新点赞计数
	go func() {
		if err := commentReplyDao.DecrementCommentLikesCount(commentId, topicId); err != nil {
			global.LOG.Errorln(err)
		}
	}()
	tx.Commit()
	return true
}

package mq

import (
	"encoding/json"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/nsqio/go-nsq"
	"log"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/core"
	dao "schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

const CommentLikeTopic = "comment_like"

type CommentLikeMessageHandler struct{}

var commentReplyDao = dao.CommentReplyDaoImpl{}

// LikeData 点赞数据
type LikeData struct {
	CommentId int64  `json:"comment_id"`
	UserId    string `json:"user_id"`
	TopicId   string `json:"topic_id"`
	Type      int    `json:"type"`
}

// CommentLikeProducer 点赞消息生产
func CommentLikeProducer(messageBody []byte) {
	err := global.NSQProducer.Publish(CommentLikeTopic, messageBody)
	if err != nil {
		global.LOG.Fatal(err)
	}
}

// CommentLikeConsumer 点赞消息消费
func CommentLikeConsumer() {
	consumer := core.InitConsumer(CommentLikeTopic)
	consumer.AddHandler(&CommentLikeMessageHandler{})
	err := consumer.ConnectToNSQD(global.CONFIG.NSQ.NsqAddr())
	if err != nil {
		log.Fatal(err)
	}
}

// HandleMessage 处理消息
func (h *CommentLikeMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var likeData LikeData
	if err := json.Unmarshal(m.Body, &likeData); err != nil {
		global.LOG.Println(err)
		return err
	}
	var err error
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err == nil {
			err = tx.Commit().Error // 确保 Commit 错误也被捕获
		} else {
			tx.Rollback()
		}
	}()

	switch likeData.Type {
	case enum.CommentLike:
		like := model.ScaCommentLikes{
			CommentId: likeData.CommentId,
			UserId:    likeData.UserId,
			TopicId:   likeData.TopicId,
		}
		if err = global.DB.Create(&like).Error; err != nil {
			tx.Rollback()
			global.LOG.Errorln(err)
			return err
		}

		err = commentReplyDao.UpdateCommentLikesCount(likeData.CommentId, likeData.TopicId)
		if err != nil {
			tx.Rollback()
			global.LOG.Errorln(err)
			return err
		}

		if err = redis.SAdd(constant.CommentLikeListRedisKey+likeData.UserId+":"+likeData.TopicId, likeData.CommentId).Err(); err != nil {
			tx.Rollback()
			return err
		}

	case enum.CommentDislike: // 取消点赞
		query, u := gplus.NewQuery[model.ScaCommentLikes]()
		query.Eq(&u.CommentId, likeData.CommentId).
			Eq(&u.UserId, likeData.UserId).
			Eq(&u.TopicId, likeData.TopicId)
		if err = gplus.Delete[model.ScaCommentLikes](query).Error; err != nil {
			tx.Rollback()
			return err
		}

		err = commentReplyDao.DecrementCommentLikesCount(likeData.CommentId, likeData.TopicId)
		if err != nil {
			tx.Rollback()
			global.LOG.Errorln(err)
			return err
		}

		if err = redis.SRem(constant.CommentLikeListRedisKey+likeData.UserId+":"+likeData.TopicId, likeData.CommentId).Err(); err != nil {
			global.LOG.Errorln(err)
			return err
		}

	default:
		global.LOG.Println("unknown comment like type")
		return nil
	}

	tx.Commit()
	return nil
}

package comment_reply_service

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

// CreateCommentReply 创建评论
func (CommentReplyService) CreateCommentReply(comment *model.ScaCommentReply) error {
	if err := global.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentListOrderByCreatedTimeDesc 通过topic_id获取评论列表
func (CommentReplyService) GetCommentListOrderByCreatedTimeDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error) {
	var comments []model.ScaCommentReply
	// 计算偏移量
	offset := (page - 1) * pageSize

	if err := global.DB.Where("topic_id =? and deleted = 0", topicID).Order("created_time desc").
		Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentListOrderByLikesDesc 通过topic_id获取评论列表
func (CommentReplyService) GetCommentListOrderByLikesDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error) {
	var comments []model.ScaCommentReply
	// 计算偏移量
	offset := (page - 1) * pageSize

	if err := global.DB.Where("topic_id =? and deleted = 0", topicID).Order("likes desc").
		Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

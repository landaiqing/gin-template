package impl

import (
	"fmt"
	"gorm.io/gorm"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

type CommentReplyDaoImpl struct{}

// CreateCommentReply 创建评论
func (CommentReplyDaoImpl) CreateCommentReply(comment *model.ScaCommentReply) error {
	if err := global.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentListOrderByCreatedTimeDesc 通过topic_id获取评论列表
func (CommentReplyDaoImpl) GetCommentListOrderByCreatedTimeDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error) {
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
func (CommentReplyDaoImpl) GetCommentListOrderByLikesDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error) {
	var comments []model.ScaCommentReply
	// 计算偏移量
	offset := (page - 1) * pageSize

	if err := global.DB.Where("topic_id =? and deleted = 0", topicID).Order("likes desc").
		Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateCommentReplyCount 更新评论
func (CommentReplyDaoImpl) UpdateCommentReplyCount(commentID int64) error {
	// 使用事务处理错误
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.ScaCommentReply{}).Where("id = ? and deleted = 0", commentID).Update("reply_count", gorm.Expr("reply_count + ?", 1))
		if result.Error != nil {
			return result.Error // 返回更新错误
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("comment not found") // 处理评论不存在的情况
		}
		return nil
	})
	return err
}

// UpdateCommentLikesCount 更新评论 likes 数量
func (CommentReplyDaoImpl) UpdateCommentLikesCount(commentID int64, topicID string) error {
	// 使用事务处理错误
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.ScaCommentReply{}).Where("id = ? and topic_id = ? and deleted = 0", commentID, topicID).Update("likes", gorm.Expr("likes + ?", 1))
		if result.Error != nil {
			return result.Error // 返回更新错误
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("comment not found") // 处理评论不存在的情况
		}
		return nil
	})
	return err
}

// DecrementCommentLikesCount 减少评论 likes 数量
func (CommentReplyDaoImpl) DecrementCommentLikesCount(commentID int64, topicID string) error {
	// 使用事务处理错误
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.ScaCommentReply{}).Where("id = ? and topic_id = ? and deleted = 0", commentID, topicID).Update("likes", gorm.Expr("likes - ?", 1))
		if result.Error != nil {
			return result.Error // 返回更新错误
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("comment not found") // 处理评论不存在的情况
		}
		return nil
	})
	return err
}

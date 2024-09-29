package dao

import "schisandra-cloud-album/model"

type CommentReplyDao interface {
	// CreateCommentReply 创建评论回复
	CreateCommentReply(comment *model.ScaCommentReply) error
	// GetCommentListOrderByCreatedTimeDesc 获取评论列表 按创建时间排序
	GetCommentListOrderByCreatedTimeDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error)
	// GetCommentListOrderByLikesDesc 获取评论列表按点赞数排序
	GetCommentListOrderByLikesDesc(topicID uint, page, pageSize int) ([]model.ScaCommentReply, error)
	// UpdateCommentReplyCount 更新评论回复数
	UpdateCommentReplyCount(commentID int64) error
	// UpdateCommentLikesCount 更新评论点赞数
	UpdateCommentLikesCount(commentID int64, topicID string) error
	// DecrementCommentLikesCount 减少评论点赞数
	DecrementCommentLikesCount(commentID int64, topicID string) error
}

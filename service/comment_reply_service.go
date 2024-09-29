package service

import (
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service/impl"
)

type CommentReplyService interface {
	// GetCommentReplyService 获取评论回复
	GetCommentReplyService(uid string, topicId string, commentId int64, pageNum int, size int) *impl.CommentResponse
	// GetCommentListService 获取评论列表
	GetCommentListService(uid string, topicId string, pageNum int, size int, isHot bool) *impl.CommentResponse
	// CommentLikeService 评论点赞
	CommentLikeService(uid string, commentId int64, topicId string) bool
	// CommentDislikeService 评论取消点赞
	CommentDislikeService(uid string, commentId int64, topicId string) bool
	// SubmitCommentService 提交评论
	SubmitCommentService(comment *model.ScaCommentReply, topicId string, uid string, images []string) bool
}

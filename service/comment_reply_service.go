package service

import "schisandra-cloud-album/model"

type CommentReplyService interface {
	CreateCommentReplyService(comment *model.ScaCommentReply) error
	UpdateCommentReplyCountService(replyId int64) error
	UpdateCommentLikesCountService(commentId int64, topicId string) error
	DecrementCommentLikesCountService(commentId int64, topicId string) error
}

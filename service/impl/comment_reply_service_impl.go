package impl

import (
	"schisandra-cloud-album/dao/impl"
	"schisandra-cloud-album/model"
)

var commentReplyDao = impl.CommentReplyDaoImpl{}

type CommentReplyServiceImpl struct{}

// CreateCommentReplyService 创建评论回复
func (CommentReplyServiceImpl) CreateCommentReplyService(comment *model.ScaCommentReply) error {
	return commentReplyDao.CreateCommentReply(comment)

}

// UpdateCommentReplyCountService 更新评论回复数
func (CommentReplyServiceImpl) UpdateCommentReplyCountService(replyId int64) error {
	return commentReplyDao.UpdateCommentReplyCount(replyId)

}

// UpdateCommentLikesCountService 更新评论点赞数
func (CommentReplyServiceImpl) UpdateCommentLikesCountService(commentId int64, topicId string) error {
	return commentReplyDao.UpdateCommentLikesCount(commentId, topicId)
}

// DecrementCommentLikesCountService 减少评论点赞数
func (CommentReplyServiceImpl) DecrementCommentLikesCountService(commentId int64, topicId string) error {
	return commentReplyDao.DecrementCommentLikesCount(commentId, topicId)
}

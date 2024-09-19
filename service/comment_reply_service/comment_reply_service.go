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

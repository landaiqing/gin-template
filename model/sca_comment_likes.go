package model

const ScaCommentLikesTableName = "sca_comment_likes"

// ScaCommentLikes 评论点赞表
type ScaCommentLikes struct {
	Id        int64  `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
	UserId    string `gorm:"column:user_id;type:varchar(20);comment:用户ID;NOT NULL" json:"user_id"`
	CommentId int64  `gorm:"column:comment_id;type:bigint(20);comment:评论ID" json:"comment_id"`
}

func (like *ScaCommentLikes) TableName() string {
	return ScaCommentLikesTableName
}

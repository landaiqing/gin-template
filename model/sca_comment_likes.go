package model

const ScaCommentLikesTableName = "sca_comment_likes"

// ScaCommentLikes 评论点赞表
type ScaCommentLikes struct {
	Id        int64  `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
	TopicId   string `gorm:"column:topic_id;type:varchar(20);NOT NULL;comment:主题ID" json:"topic_id"`
	UserId    string `gorm:"column:user_id;type:varchar(20);comment:用户ID;NOT NULL" json:"user_id"`
	CommentId int64  `gorm:"column:comment_id;type:bigint(20);comment:评论ID;NOT NULL" json:"comment_id"`
	Type      int    `gorm:"column:type;type:int(11);comment:类型，0 点赞 1 踩;NOT NULL" json:"type"`
}

func (like *ScaCommentLikes) TableName() string {
	return ScaCommentLikesTableName
}

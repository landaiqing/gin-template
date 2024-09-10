package model

import (
	"gorm.io/gorm"
	"time"
)

// ScaCommentReply 评论表
type ScaCommentReply struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	UserId      string    `gorm:"column:user_id;type:varchar(20);comment:评论用户id" json:"user_id"`
	TopicId     string    `gorm:"column:topic_id;type:varchar(20);comment:评论话题id" json:"topic_id"`
	TopicType   int       `gorm:"column:topic_type;type:int(11);comment:话题类型" json:"topic_type"`
	Content     string    `gorm:"column:content;type:longtext;comment:评论内容" json:"content"`
	CommentType int       `gorm:"column:comment_type;type:int(11);comment:评论类型 0评论 1 回复" json:"comment_type"`
	ReplyId     string    `gorm:"column:reply_id;type:varchar(20);comment:回复目标id" json:"reply_id"`
	ReplyUser   string    `gorm:"column:reply_user;type:varchar(20);comment:回复人id" json:"reply_user"`
	Username    string    `gorm:"column:username;type:varchar(20);comment:评论用户昵称" json:"username"`
	Avatar      string    `gorm:"column:avatar;type:longtext;comment:评论用户头像" json:"avatar"`
	Author      int       `gorm:"column:author;type:int(11);comment:评论回复是否作者  0否 1是" json:"author"`
	Likes       int64     `gorm:"column:likes;type:bigint(20);comment:点赞数" json:"likes"`
	ReplyCount  int64     `gorm:"column:reply_count;type:bigint(20);comment:回复数量" json:"reply_count"`
	PicUrls     string    `gorm:"column:pic_urls;type:longtext;comment:图片链接" json:"pic_urls"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`
	Deleted     int       `gorm:"column:deleted;type:int(11);default:0;comment:是否删除 0未删除 1 已删除" json:"deleted"`
	CreatedBy   string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"created_by"`
	UpdateBy    string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"update_by"`
}

func (comment *ScaCommentReply) TableName() string {
	return "sca_comment_reply"
}

// BeforeCreate 创建前回调
func (comment *ScaCommentReply) BeforeCreate(tx *gorm.DB) (err error) {
	comment.CreatedBy = comment.UserId
	return nil
}

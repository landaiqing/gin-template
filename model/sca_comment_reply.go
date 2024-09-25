package model

import (
	"gorm.io/gorm"
	"time"
)

const ScaCommentReplyTable = "sca_comment_reply"

// ScaCommentReply 评论表
type ScaCommentReply struct {
	Id              int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	UserId          string    `gorm:"column:user_id;type:varchar(20);comment:评论用户id" json:"user_id"`
	TopicId         string    `gorm:"column:topic_id;type:varchar(20);comment:评论话题id" json:"topic_id"`
	TopicType       int       `gorm:"column:topic_type;type:int(11);comment:话题类型" json:"-"`
	Content         string    `gorm:"column:content;type:longtext;comment:评论内容" json:"content"`
	CommentType     int       `gorm:"column:comment_type;type:int(11);comment:评论类型 0评论 1 回复" json:"-"`
	ReplyTo         int64     `gorm:"column:reply_to;type:bigint(20);comment:回复子评论id" json:"reply_to"`
	ReplyId         int64     `gorm:"column:reply_id;type:varchar(20);comment:回复目标id" json:"reply_id,omitempty"`
	ReplyUser       string    `gorm:"column:reply_user;type:varchar(20);comment:回复人id" json:"reply_user,omitempty"`
	Author          int       `gorm:"column:author;type:int(11);default:0;comment:评论回复是否作者  0否 1是" json:"author"`
	Likes           int64     `gorm:"column:likes;type:bigint(20);default:0;comment:点赞数" json:"likes"`
	ReplyCount      int64     `gorm:"column:reply_count;type:bigint(20);default:0;comment:回复数量" json:"reply_count"`
	CreatedTime     time.Time `gorm:"column:created_time;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"-"`
	Deleted         int       `gorm:"column:deleted;type:int(11);default:0;comment:是否删除 0未删除 1 已删除" json:"-"`
	CreatedBy       string    `gorm:"column:created_by;type:varchar(32);comment:创建人" json:"-"`
	UpdateBy        string    `gorm:"column:update_by;type:varchar(32);comment:更新人" json:"-"`
	CommentIp       string    `gorm:"column:comment_ip;type:varchar(20);comment:评论ip" json:"-"`
	Location        string    `gorm:"column:location;type:varchar(20);comment:评论地址" json:"location"`
	Browser         string    `gorm:"column:browser;type:varchar(20);comment:评论浏览器" json:"browser"`
	OperatingSystem string    `gorm:"column:operating_system;type:varchar(20);comment:评论操作系统" json:"operating_system"`
	CommentOrder    int64     `gorm:"column:comment_order;type:bigint(20);default:0;comment:评论排序" json:"-"`
	Agent           string    `gorm:"column:agent;type:varchar(255);comment:客户端设备信息" json:"agent"`
}

func (comment *ScaCommentReply) TableName() string {
	return ScaCommentReplyTable
}

// BeforeCreate 创建前回调
func (comment *ScaCommentReply) BeforeCreate(tx *gorm.DB) (err error) {
	comment.CreatedBy = comment.UserId
	return nil
}

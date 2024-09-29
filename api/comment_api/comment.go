package comment_api

import (
	"schisandra-cloud-album/service"
	"sync"
	"time"
)

type CommentAPI struct{}

var wg sync.WaitGroup
var mx sync.Mutex
var commentReplyService = service.Service.CommentReplyService

// CommentImages 评论图片
type CommentImages struct {
	TopicId   string   `json:"topic_id" bson:"topic_id" required:"true"`
	CommentId int64    `json:"comment_id" bson:"comment_id" required:"true"`
	UserId    string   `json:"user_id" bson:"user_id" required:"true"`
	Images    [][]byte `json:"images" bson:"images" required:"true"`
	CreatedAt string   `json:"created_at" bson:"created_at" required:"true"`
}

// CommentContent 评论内容
type CommentContent struct {
	NickName        string    `json:"nickname"`
	Avatar          string    `json:"avatar"`
	Level           int       `json:"level,omitempty"`
	Id              int64     `json:"id"`
	UserId          string    `json:"user_id"`
	TopicId         string    `json:"topic_id"`
	Content         string    `json:"content"`
	ReplyTo         int64     `json:"reply_to,omitempty"`
	ReplyId         int64     `json:"reply_id,omitempty"`
	ReplyUser       string    `json:"reply_user,omitempty"`
	ReplyUsername   string    `json:"reply_username,omitempty"`
	Author          int       `json:"author"`
	Likes           int64     `json:"likes"`
	ReplyCount      int64     `json:"reply_count"`
	CreatedTime     time.Time `json:"created_time"`
	Location        string    `json:"location"`
	Browser         string    `json:"browser"`
	OperatingSystem string    `json:"operating_system"`
	IsLiked         bool      `json:"is_liked" default:"false"`
	Images          []string  `json:"images,omitempty"`
}

// CommentResponse 评论返回值
type CommentResponse struct {
	Size     int              `json:"size"`
	Total    int64            `json:"total"`
	Current  int              `json:"current"`
	Comments []CommentContent `json:"comments"`
}

var likeChannel = make(chan CommentLikeRequest, 1000)
var cancelLikeChannel = make(chan CommentLikeRequest, 1000)

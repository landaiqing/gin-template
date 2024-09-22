package comment_api

import (
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
)

type CommentAPI struct{}

var commentReplyService = service.Service.CommentReplyService

type CommentImages struct {
	TopicId   string   `json:"topic_id" bson:"topic_id" required:"true"`
	CommentId int64    `json:"comment_id" bson:"comment_id" required:"true"`
	UserId    string   `json:"user_id" bson:"user_id" required:"true"`
	Images    []string `json:"image_url" bson:"images" required:"true"`
	CreatedAt string   `json:"created_at" bson:"created_at" required:"true"`
}

type CommentData struct {
	Comment model.ScaCommentReply `json:"comment"`
	Images  []string              `json:"images,omitempty"`
}

type CommentResponse struct {
	Size     int           `json:"size"`
	Total    int64         `json:"total"`
	Current  int           `json:"current"`
	Comments []CommentData `json:"comments"`
}

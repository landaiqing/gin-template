package dto

type CommentRequest struct {
	Content string   `json:"content"`
	Images  []string `json:"images"`
	UserID  string   `json:"user_id"`
	TopicId string   `json:"topic_id"`
	Author  string   `json:"author"`
}
type ReplyCommentRequest struct {
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	UserID    string   `json:"user_id"`
	TopicId   string   `json:"topic_id"`
	ReplyId   string   `json:"reply_id"`
	ReplyUser string   `json:"reply_user"`
	Author    string   `json:"author"`
}
type CommentListRequest struct {
	TopicId string `json:"topic_id"`
	Page    int    `json:"page" default:"1"`
	Size    int    `json:"size" default:"10"`
}

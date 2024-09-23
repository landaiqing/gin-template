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
	ReplyId   int64    `json:"reply_id"`
	ReplyUser string   `json:"reply_user"`
	Author    string   `json:"author"`
}

type ReplyReplyRequest struct {
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	UserID    string   `json:"user_id"`
	TopicId   string   `json:"topic_id"`
	ReplyTo   int64    `json:"reply_to"`
	ReplyId   int64    `json:"reply_id"`
	ReplyUser string   `json:"reply_user"`
	Author    string   `json:"author"`
}

type CommentListRequest struct {
	TopicId string `json:"topic_id"`
	Page    int    `json:"page" default:"1"`
	Size    int    `json:"size" default:"5"`
}
type ReplyListRequest struct {
	TopicId   string `json:"topic_id"`
	CommentId int64  `json:"comment_id"`
	Page      int    `json:"page" default:"1"`
	Size      int    `json:"size" default:"5"`
}

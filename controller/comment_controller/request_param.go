package comment_controller

// CommentRequest 评论请求参数
type CommentRequest struct {
	Content string   `json:"content" binding:"required"`
	Images  []string `json:"images"`
	UserID  string   `json:"user_id" binding:"required"`
	TopicId string   `json:"topic_id" binding:"required"`
	Author  string   `json:"author" binding:"required"`
	Key     string   `json:"key" binding:"required"`
	Point   []int64  `json:"point" binding:"required"`
}

// ReplyCommentRequest 回复评论请求参数
type ReplyCommentRequest struct {
	Content   string   `json:"content" binding:"required"`
	Images    []string `json:"images"`
	UserID    string   `json:"user_id" binding:"required"`
	TopicId   string   `json:"topic_id" binding:"required"`
	ReplyId   int64    `json:"reply_id" binding:"required"`
	ReplyUser string   `json:"reply_user" binding:"required"`
	Author    string   `json:"author" binding:"required"`
	Key       string   `json:"key" binding:"required"`
	Point     []int64  `json:"point" binding:"required"`
}

// ReplyReplyRequest 回复回复请求参数
type ReplyReplyRequest struct {
	Content   string   `json:"content" binding:"required"`
	Images    []string `json:"images"`
	UserID    string   `json:"user_id" binding:"required"`
	TopicId   string   `json:"topic_id" binding:"required"`
	ReplyTo   int64    `json:"reply_to" binding:"required"`
	ReplyId   int64    `json:"reply_id" binding:"required"`
	ReplyUser string   `json:"reply_user" binding:"required"`
	Author    string   `json:"author" binding:"required"`
	Key       string   `json:"key" binding:"required"`
	Point     []int64  `json:"point" binding:"required"`
}

// CommentListRequest 评论列表请求参数
type CommentListRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	TopicId string `json:"topic_id" binding:"required"`
	Page    int    `json:"page" default:"1"`
	Size    int    `json:"size" default:"5"`
	IsHot   bool   `json:"is_hot" default:"true"`
}

// ReplyListRequest 回复列表请求参数
type ReplyListRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	TopicId   string `json:"topic_id" binding:"required"`
	CommentId int64  `json:"comment_id" binding:"required"`
	Page      int    `json:"page" default:"1"`
	Size      int    `json:"size" default:"5"`
}

// CommentLikeRequest 点赞评论的请求参数
type CommentLikeRequest struct {
	TopicId   string `json:"topic_id" binding:"required"`
	CommentId int64  `json:"comment_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

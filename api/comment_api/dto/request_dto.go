package dto

type CommentRequest struct {
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

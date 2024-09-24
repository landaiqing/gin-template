package comment_api

import (
	"encoding/base64"
	"errors"
	"io"
	"regexp"
	"schisandra-cloud-album/service"
	"strings"
	"sync"
	"time"
)

type CommentAPI struct{}

var wg sync.WaitGroup
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
	Dislikes        int64     `json:"dislikes"`
	Location        string    `json:"location"`
	Browser         string    `json:"browser"`
	OperatingSystem string    `json:"operating_system"`
	Images          []string  `json:"images,omitempty"`
}

// CommentResponse 评论返回值
type CommentResponse struct {
	Size     int              `json:"size"`
	Total    int64            `json:"total"`
	Current  int              `json:"current"`
	Comments []CommentContent `json:"comments"`
}

// base64ToBytes 将base64字符串转换为字节数组
func base64ToBytes(base64Str string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str))
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("failed to decode base64 string")
	}
	return data, nil
}

// processImages 处理图片，将 base64 字符串转换为字节数组
func processImages(images []string) ([][]byte, error) {
	var imagesData [][]byte
	dataChan := make(chan []byte, len(images)) // 创建一个带缓冲的 channel
	re := regexp.MustCompile(`^data:image/\w+;base64,`)

	for _, img := range images {
		wg.Add(1) // 增加 WaitGroup 的计数
		go func(img string) {
			defer wg.Done() // 函数结束时减少计数

			imgWithoutPrefix := re.ReplaceAllString(img, "")
			data, err := base64ToBytes(imgWithoutPrefix)
			if err != nil {
				return // 出错时直接返回
			}
			dataChan <- data // 将结果发送到 channel
		}(img)
	}

	wg.Wait()       // 等待所有 goroutine 完成
	close(dataChan) // 关闭 channel

	for data := range dataChan { // 收集所有结果
		imagesData = append(imagesData, data)
	}

	return imagesData, nil
}

// getMimeType 获取 MIME 类型
func getMimeType(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream" // 默认类型
	}

	// 判断 JPEG
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// 判断 PNG
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// 判断 GIF
	if len(data) >= 6 && data[0] == 'G' && data[1] == 'I' && data[2] == 'F' {
		return "image/gif"
	}
	// 判断 WEBP
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return "image/webp"
	}
	// 判断svg
	if len(data) >= 4 && data[0] == '<' && data[1] == '?' && data[2] == 'x' && data[3] == 'm' {
		return "image/svg+xml"
	}
	// 判断JPG
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}

	return "application/octet-stream" // 默认类型
}

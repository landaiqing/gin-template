package utils

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"
)

// 获取网络图片
var client = &http.Client{
	Timeout: 5 * time.Second, // 超时时间
}

// GenerateAvatar 用于生成用户头像
func GenerateAvatar(userId string) (baseImg string, err error) {

	path := "https://api.multiavatar.com/" + userId + ".png"

	// 创建请求
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", errors.New("image request error")
	}

	// 发送请求并获取响应
	respImg, err := client.Do(request)
	if err != nil {
		return "", errors.New("failed to fetch image")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(respImg.Body)

	// 读取图片数据
	imgByte, err := io.ReadAll(respImg.Body)
	if err != nil {
		return "", errors.New("failed to read image data")
	}

	// 判断文件类型，生成一个前缀
	mimeType := http.DetectContentType(imgByte)
	switch mimeType {
	case "image/png":
		baseImg = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgByte)
	default:
		return "", errors.New("unsupported image type")
	}
	return baseImg, nil
}

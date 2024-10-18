package utils

import (
	"encoding/base64"
	"io"
	"net/http"
	"schisandra-cloud-album/global"
	"strconv"
	"time"
)

const (
	numParts = 4 // 分成4块
)

func GenerateAvatar(userId string) (baseImg string) {
	path := "https://api.multiavatar.com/" + userId + ".png"

	// 创建 HTTP 客户端并设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发送 HEAD 请求获取图片大小
	headReq, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		global.LOG.Error(err)
		return ""
	}

	respHead, err := client.Do(headReq)
	if err != nil {
		global.LOG.Error(err)
		return ""
	}
	defer respHead.Body.Close()

	// 获取图片大小
	contentLength := respHead.ContentLength
	if contentLength <= 0 {
		return ""
	}

	imgChunks := make([][]byte, numParts) // 存储每个部分的图片数据

	// 启动多个 goroutine 下载分块
	for i := 0; i < numParts; i++ {
		wg.Add(1)
		go func(part int) {
			defer wg.Done()
			start := (contentLength / int64(numParts)) * int64(part)
			end := start + (contentLength / int64(numParts)) - 1
			if part == numParts-1 {
				end = contentLength - 1 // 最后一部分下载到文件末尾
			}

			// 创建 RANGE 请求
			rangeHeader := "bytes=" + strconv.FormatInt(start, 10) + "-" + strconv.FormatInt(end, 10)
			request, err := http.NewRequest("GET", path, nil)
			if err != nil {
				global.LOG.Error(err)
				return
			}
			request.Header.Set("Range", rangeHeader)

			respImg, err := client.Do(request)
			if err != nil {
				global.LOG.Error(err)
				return
			}
			defer respImg.Body.Close()

			// 读取图片数据
			imgByte, err := io.ReadAll(respImg.Body)
			if err != nil {
				global.LOG.Error(err)
				return
			}

			// 存储分块
			imgChunks[part] = imgByte
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成

	// 合并所有部分
	var fullImg []byte
	for _, chunk := range imgChunks {
		fullImg = append(fullImg, chunk...)
	}

	// 判断文件类型，生成一个前缀
	mimeType := http.DetectContentType(fullImg)
	switch mimeType {
	case "image/png":
		baseImg = "data:image/png;base64," + base64.StdEncoding.EncodeToString(fullImg)
	default:
		return ""
	}

	return baseImg
}

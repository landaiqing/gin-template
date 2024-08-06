package api

import (
	"schisandra-cloud-album/api/auth_api"
	"schisandra-cloud-album/api/captcha_api"
)

// Apis 统一导出的api
type Apis struct {
	AuthApi    auth_api.AuthAPI
	CaptchaAPI captcha_api.CaptchaAPI
}

// Api new函数实例化，实例化完成后会返回结构体地指针类型
var Api = new(Apis)

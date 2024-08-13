package api

import (
	"schisandra-cloud-album/api/captcha_api"
	"schisandra-cloud-album/api/sms_api"
	"schisandra-cloud-album/api/user_api"
)

// Apis 统一导出的api
type Apis struct {
	UserApi    user_api.UserAPI
	CaptchaApi captcha_api.CaptchaAPI
	SmsApi     sms_api.SmsAPI
}

// Api new函数实例化，实例化完成后会返回结构体地指针类型
var Api = new(Apis)

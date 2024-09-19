package api

import (
	"schisandra-cloud-album/api/captcha_api"
	"schisandra-cloud-album/api/client_api"
	"schisandra-cloud-album/api/comment_api"
	"schisandra-cloud-album/api/oauth_api"
	"schisandra-cloud-album/api/permission_api"
	"schisandra-cloud-album/api/role_api"
	"schisandra-cloud-album/api/sms_api"
	"schisandra-cloud-album/api/user_api"
	"schisandra-cloud-album/api/websocket_api"
)

// Apis 统一导出的api
type Apis struct {
	UserApi       user_api.UserAPI
	CaptchaApi    captcha_api.CaptchaAPI
	SmsApi        sms_api.SmsAPI
	OAuthApi      oauth_api.OAuthAPI
	WebsocketApi  websocket_api.WebsocketAPI
	RoleApi       role_api.RoleAPI
	PermissionApi permission_api.PermissionAPI
	ClientApi     client_api.ClientAPI
	CommonApi     comment_api.CommentAPI
}

// Api new函数实例化，实例化完成后会返回结构体地指针类型
var Api = new(Apis)

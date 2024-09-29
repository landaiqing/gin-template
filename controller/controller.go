package controller

import (
	"schisandra-cloud-album/controller/captcha_controller"
	"schisandra-cloud-album/controller/client_controller"
	"schisandra-cloud-album/controller/comment_controller"
	"schisandra-cloud-album/controller/oauth_controller"
	"schisandra-cloud-album/controller/permission_controller"
	"schisandra-cloud-album/controller/role_controller"
	"schisandra-cloud-album/controller/sms_controller"
	"schisandra-cloud-album/controller/user_controller"
	"schisandra-cloud-album/controller/websocket_controller"
)

// Controllers 统一导出的控制器接口
type Controllers struct {
	UserController       user_controller.UserController
	CaptchaController    captcha_controller.CaptchaController
	SmsController        sms_controller.SmsController
	OAuthController      oauth_controller.OAuthController
	WebsocketController  websocket_controller.WebsocketController
	RoleController       role_controller.RoleController
	PermissionController permission_controller.PermissionController
	ClientController     client_controller.ClientController
	CommonController     comment_controller.CommentController
}

// Controller new函数实例化，实例化完成后会返回结构体地指针类型
var Controller = new(Controllers)

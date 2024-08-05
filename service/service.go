package service

import (
	"schisandra-cloud-album/service/auth_service"
)

// Services 统一导出的service
type Services struct {
	AuthService auth_service.AuthService
}

// Service new函数实例化，实例化完成后会返回结构体地指针类型
var Service = new(Services)

package user_controller

import (
	"schisandra-cloud-album/service/impl"
	"sync"
)

type UserController struct{}

var mu sync.Mutex
var userService = impl.UserServiceImpl{}
var userDeviceService = impl.UserDeviceServiceImpl{}

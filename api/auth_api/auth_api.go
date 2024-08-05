package auth_api

import (
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/service"
)

var authService = service.Service.AuthService

// GetUserList
// @Summary 获取所有用户列表
// @Tags 鉴权模块
// @Success 200 {string} json
// @Router /api/auth/user/List [get]
func (AuthAPI) GetUserList(c *gin.Context) {
	userList := authService.GetUserList()
	result.OkWithData(userList, c)
}

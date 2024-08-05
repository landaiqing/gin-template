package auth_api

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/model"
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

// QueryUserByUsername
// @Summary 根据用户名查询用户
// @Tags 鉴权模块
// @Param username query string true "用户名"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_username [get]
func (AuthAPI) QueryUserByUsername(c *gin.Context) {
	username := c.Query("username")
	user := authService.QueryUserByName(username)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage("用户不存在！", c)
		return
	}
	result.OkWithData(user, c)
}

// QueryUserByUuid
// @Summary 根据uuid查询用户
// @Tags 鉴权模块
// @Param uuid query string true "用户uuid"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_uuid [get]
func (AuthAPI) QueryUserByUuid(c *gin.Context) {
	uuid := c.Query("uuid")
	user := authService.QueryUserByUuid(uuid)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage("用户不存在！", c)
		return
	}
	result.OkWithData(user, c)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 鉴权模块
// @Param uuid query string true "用户uuid"
// @Success 200 {string} json
// @Router /api/auth/user/delete [delete]
func (AuthAPI) DeleteUser(c *gin.Context) {
	uuid := c.Query("uuid")
	err := authService.DeleteUser(uuid)
	if err != nil {
		result.FailWithMessage("用户删除失败！", c)
		return
	}
	result.OkWithMessage("用户删除成功！", c)
}

// QueryUserByPhone 根据手机号查询用户
// @Summary 根据手机号查询用户
// @Tags 鉴权模块
// @Param phone query string true "手机号"
// @Success 200 {string} json
// @Router /api/auth/user/query_by_phone [get]
func (AuthAPI) QueryUserByPhone(c *gin.Context) {
	phone := c.Query("phone")
	user := authService.QueryUserByPhone(phone)
	if reflect.DeepEqual(user, model.ScaAuthUser{}) {
		result.FailWithMessage("用户不存在！", c)
		return
	}
	result.OkWithData(user, c)
}

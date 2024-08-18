package oauth_api

import (
	"encoding/json"
	"schisandra-cloud-album/api/user_api/dto"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
	"schisandra-cloud-album/utils"
	"time"
)

type OAuthAPI struct{}

var userService = service.Service.UserService
var userRoleService = service.Service.UserRoleService
var userSocialService = service.Service.UserSocialService
var rolePermissionService = service.Service.RolePermissionService
var permissionServiceService = service.Service.PermissionService
var roleService = service.Service.RoleService

var script = `
        <script>
        window.opener.postMessage('%s', '%s');
        window.close();
        </script>
        `

// HandelUserLogin 处理用户登录
func HandelUserLogin(user model.ScaAuthUser) (bool, map[string]interface{}) {
	ids, err := userRoleService.GetUserRoleIdsByUserId(user.ID)
	if err != nil {
		return false, nil
	}
	permissionIds := rolePermissionService.QueryPermissionIdsByRoleId(ids)
	permissions, err := permissionServiceService.GetPermissionsByIds(permissionIds)
	if err != nil {
		return false, nil
	}
	serializedPermissions, err := json.Marshal(permissions)
	if err != nil {
		return false, nil
	}
	wrong := redis.Set(constant.UserAuthPermissionRedisKey+*user.UID, serializedPermissions, 0).Err()
	if wrong != nil {
		return false, nil
	}
	roleList, err := roleService.GetRoleListByIds(ids)
	if err != nil {
		return false, nil
	}
	serializedRoleList, err := json.Marshal(roleList)
	if err != nil {
		return false, nil
	}
	er := redis.Set(constant.UserAuthRoleRedisKey+*user.UID, serializedRoleList, 0).Err()
	if er != nil {
		return false, nil
	}
	accessToken, refreshToken, expiresAt := utils.GenerateAccessTokenAndRefreshToken(utils.JWTPayload{UserID: user.UID, RoleID: ids})

	data := dto.ResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UID:          user.UID,
	}
	fail := redis.Set(constant.UserLoginTokenRedisKey+*user.UID, data, time.Hour*24*7).Err()
	if fail != nil {
		return false, nil
	}
	responseData := map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
		"success": true,
	}
	return true, responseData
}

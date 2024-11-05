package oauth_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/common/types"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/service/impl"
	"schisandra-cloud-album/utils"
)

type OAuthController struct{}

var userSocialService = impl.UserSocialServiceImpl{}
var userService = impl.UserServiceImpl{}

type Token struct {
	AccessToken string `json:"access_token"`
}

var script = `
        <script>
        window.opener.postMessage('%s', '%s');
        window.close();
        </script>
        `

func HandleLoginResponse(c *gin.Context, uid string) {
	// 查询用户信息
	user := userService.QueryUserByUuidService(&uid)

	var accessToken, refreshToken string
	var err error
	var wg sync.WaitGroup
	var accessTokenErr error

	wg.Add(2) // 增加计数器，等待两个协程完成

	// 使用goroutine生成accessToken
	go func() {
		defer wg.Done() // 完成时减少计数器
		accessToken, accessTokenErr = utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: &uid})
	}()

	// 使用goroutine生成refreshToken
	go func() {
		defer wg.Done() // 完成时减少计数器
		refreshToken = utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: &uid}, time.Hour*24*7)
	}()

	// 等待两个协程完成
	wg.Wait()

	// 检查生成accessToken时是否有错误
	if accessTokenErr != nil {
		global.LOG.Error(accessTokenErr)
		return
	}

	data := types.ResponseData{
		AccessToken: accessToken,
		UID:         &uid,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
	}
	// 设置session
	sessionData := utils.SessionData{
		RefreshToken: refreshToken,
		UID:          uid,
	}
	if err = utils.SetSession(c, constant.SessionKey, sessionData); err != nil {
		return
	}
	redisTokenData := types.RedisToken{
		AccessToken: accessToken,
		UID:         uid,
	}
	// 将数据存入redis
	if err = redis.Set(constant.UserLoginTokenRedisKey+uid, redisTokenData, time.Hour*24*7).Err(); err != nil {
		global.LOG.Error(err)
		return
	}

	responseData := result.Response{
		Data:    data,
		Message: "success",
		Code:    200,
		Success: true,
	}

	tokenData, err := json.Marshal(responseData)
	if err != nil {
		global.LOG.Error(err)
		return
	}

	formattedScript := fmt.Sprintf(script, tokenData, global.CONFIG.System.WebURL())
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedScript))
	return
}

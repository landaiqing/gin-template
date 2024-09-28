package oauth_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"gorm.io/gorm"
	"net/http"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"schisandra-cloud-album/service"
	"schisandra-cloud-album/utils"
	"sync"
	"time"
)

var mu sync.Mutex

type OAuthAPI struct{}

var userService = service.Service.UserService
var userSocialService = service.Service.UserSocialService
var userDeviceService = service.Service.UserDeviceService

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
	res, data := HandelUserLogin(uid, c)
	if !res {
		return
	}

	tokenData, err := json.Marshal(data)
	if err != nil {
		global.LOG.Error(err)
		return
	}

	formattedScript := fmt.Sprintf(script, tokenData, global.CONFIG.System.WebURL())
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedScript))
	return
}

// HandelUserLogin 处理用户登录
func HandelUserLogin(userId string, c *gin.Context) (bool, result.Response) {
	// 使用goroutine生成accessToken
	accessTokenChan := make(chan string)
	errChan := make(chan error)
	go func() {
		accessToken, err := utils.GenerateAccessToken(utils.AccessJWTPayload{UserID: &userId})
		if err != nil {
			errChan <- err
			return
		}
		accessTokenChan <- accessToken
	}()

	// 使用goroutine生成refreshToken
	refreshTokenChan := make(chan string)
	expiresAtChan := make(chan int64)
	go func() {
		refreshToken, expiresAt := utils.GenerateRefreshToken(utils.RefreshJWTPayload{UserID: &userId}, time.Hour*24*7)
		refreshTokenChan <- refreshToken
		expiresAtChan <- expiresAt
	}()

	// 等待accessToken和refreshToken生成完成
	var accessToken string
	var refreshToken string
	var expiresAt int64
	var err error
	select {
	case accessToken = <-accessTokenChan:
	case err = <-errChan:
		global.LOG.Error(err)
		return false, result.Response{}
	}
	select {
	case refreshToken = <-refreshTokenChan:
	case expiresAt = <-expiresAtChan:
	}

	data := ResponseData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UID:          &userId,
	}
	wrong := utils.SetSession(c, constant.SessionKey, data)
	if wrong != nil {
		return false, result.Response{}
	}
	// 使用goroutine将数据存入redis
	redisErrChan := make(chan error)
	go func() {
		fail := redis.Set(constant.UserLoginTokenRedisKey+userId, data, time.Hour*24*7).Err()
		if fail != nil {
			redisErrChan <- fail
			return
		}
		redisErrChan <- nil
	}()

	// 等待redis操作完成
	redisErr := <-redisErrChan
	if redisErr != nil {
		global.LOG.Error(redisErr)
		return false, result.Response{}
	}
	responseData := result.Response{
		Data:    data,
		Message: "success",
		Code:    200,
		Success: true,
	}
	return true, responseData
}

// GetUserLoginDevice 获取用户登录设备
func (OAuthAPI) GetUserLoginDevice(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		return
	}
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		global.LOG.Errorln("user-agent is empty")
		return
	}
	ua := useragent.New(userAgent)

	ip := utils.GetClientIP(c)
	location, err := global.IP2Location.SearchByStr(ip)
	location = utils.RemoveZeroAndAdjust(location)
	if err != nil {
		global.LOG.Errorln(err)
		return
	}
	isBot := ua.Bot()
	browser, browserVersion := ua.Browser()
	os := ua.OS()
	mobile := ua.Mobile()
	mozilla := ua.Mozilla()
	platform := ua.Platform()
	engine, engineVersion := ua.Engine()
	device := model.ScaAuthUserDevice{
		UserID:          &userId,
		IP:              &ip,
		Location:        &location,
		Agent:           userAgent,
		Browser:         &browser,
		BrowserVersion:  &browserVersion,
		OperatingSystem: &os,
		Mobile:          &mobile,
		Bot:             &isBot,
		Mozilla:         &mozilla,
		Platform:        &platform,
		EngineName:      &engine,
		EngineVersion:   &engineVersion,
	}
	mu.Lock()
	defer mu.Unlock()
	userDevice, err := userDeviceService.GetUserDeviceByUIDIPAgent(userId, ip, userAgent)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = userDeviceService.AddUserDevice(&device)
		if err != nil {
			global.LOG.Errorln(err)
			return
		}
		return
	} else {
		err := userDeviceService.UpdateUserDevice(userDevice.ID, &device)
		if err != nil {
			global.LOG.Errorln(err)
			return
		}
		return
	}
}

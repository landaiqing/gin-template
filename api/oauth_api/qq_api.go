package oauth_api

import (
	"encoding/json"
	"errors"
	"fmt"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
	"net/http"
	"schisandra-cloud-album/common/enum"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
	"strconv"
)

type AuthQQme struct {
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid"`
}
type QQToken struct {
	AccessToken  string `json:"access_token"`
	ExpireIn     string `json:"expire_in"`
	RefreshToken string `json:"refresh_token"`
}

type QQUserInfo struct {
	City            string `json:"city"`
	Figureurl       string `json:"figureurl"`
	Figureurl1      string `json:"figureurl_1"`
	Figureurl2      string `json:"figureurl_2"`
	FigureurlQq     string `json:"figureurl_qq"`
	FigureurlQq1    string `json:"figureurl_qq_1"`
	FigureurlQq2    string `json:"figureurl_qq_2"`
	Gender          string `json:"gender"`
	GenderType      int    `json:"gender_type"`
	IsLost          int    `json:"is_lost"`
	IsYellowVip     string `json:"is_yellow_vip"`
	IsYellowYearVip string `json:"is_yellow_year_vip"`
	Level           string `json:"level"`
	Msg             string `json:"msg"`
	Nickname        string `json:"nickname"`
	Province        string `json:"province"`
	Ret             int    `json:"ret"`
	Vip             string `json:"vip"`
	Year            string `json:"year"`
	YellowVipLevel  string `json:"yellow_vip_level"`
}

// GetQQRedirectUrl 获取登录地址
// @Summary 获取QQ登录地址
// @Description 获取QQ登录地址
// @Tags QQ OAuth
// @Produce  json
// @Success 200 {string} string "登录地址"
// @Router /api/oauth/qq/get_url [get]
func (OAuthAPI) GetQQRedirectUrl(c *gin.Context) {
	state := c.Query("state")
	clientId := global.CONFIG.OAuth.QQ.ClientID
	redirectURI := global.CONFIG.OAuth.QQ.RedirectURI
	url := "https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=" + clientId + "&redirect_uri=" + redirectURI + "&state=" + state
	result.OkWithData(url, c)
	return
}

// GetQQTokenAuthUrl 通过code获取token认证url
func GetQQTokenAuthUrl(code string) string {
	clientId := global.CONFIG.OAuth.QQ.ClientID
	clientSecret := global.CONFIG.OAuth.QQ.ClientSecret
	redirectURI := global.CONFIG.OAuth.QQ.RedirectURI
	return fmt.Sprintf(
		"https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&fmt=json",
		clientId, clientSecret, code, redirectURI,
	)
}

// GetQQToken 获取 token
func GetQQToken(url string) (*QQToken, error) {

	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		global.LOG.Error(err)
		return nil, err
	}

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	//将响应体解析为 token，并返回
	var token QQToken
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	return &token, nil
}

// GetQQUserOpenID 获取用户 openid
func GetQQUserOpenID(token *QQToken) (*AuthQQme, error) {

	// 形成请求
	var userInfoUrl = "https://graph.qq.com/oauth2.0/me?access_token=" + token.AccessToken + "&fmt=json"
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		global.LOG.Error(err)
		return nil, err
	}

	// 将响应体解析为 AuthQQme，并返回
	var authQQme AuthQQme
	if err = json.NewDecoder(res.Body).Decode(&authQQme); err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	return &authQQme, nil
}

// GetQQUserUserInfo 获取用户信息
func GetQQUserUserInfo(token *QQToken, openId string) (map[string]interface{}, error) {

	clientId := global.CONFIG.OAuth.QQ.ClientID
	// 形成请求
	var userInfoUrl = "https://graph.qq.com/user/get_user_info?access_token=" + token.AccessToken + "&oauth_consumer_key=" + clientId + "&openid=" + openId
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		global.LOG.Error(err)
		return nil, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	return userInfo, nil
}

// QQCallback QQ登录回调
// @Summary QQ登录回调
// @Description QQ登录回调
// @Tags QQ OAuth
// @Produce  json
// @Router /api/oauth/qq/callback [get]
func (OAuthAPI) QQCallback(c *gin.Context) {
	var err error
	// 获取 code
	var code = c.Query("code")
	if code == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	// 通过 code, 获取 token
	var tokenAuthUrl = GetQQTokenAuthUrl(code)
	var token *QQToken
	if token, err = GetQQToken(tokenAuthUrl); err != nil {
		global.LOG.Error(err)
		return
	}
	// 通过 token，获取 openid
	authQQme, err := GetQQUserOpenID(token)
	if err != nil {
		global.LOG.Error(err)
		return
	}

	// 通过token，获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = GetQQUserUserInfo(token, authQQme.OpenID); err != nil {
		global.LOG.Error(err)
		return
	}

	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	var qqUserInfo QQUserInfo
	err = json.Unmarshal(userInfoBytes, &qqUserInfo)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	userSocial, err := userSocialService.QueryUserSocialByOpenID(authQQme.OpenID, enum.OAuthSourceQQ)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 第一次登录，创建用户
		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		location := qqUserInfo.Province + "|" + qqUserInfo.City
		user := model.ScaAuthUser{
			UID:      &uidStr,
			Username: &authQQme.OpenID,
			Nickname: &qqUserInfo.Nickname,
			Avatar:   &qqUserInfo.FigureurlQq1,
			Gender:   &qqUserInfo.Gender,
			Location: &location,
		}
		addUser, err := userService.AddUser(user)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		qq := enum.OAuthSourceQQ
		userSocial = model.ScaAuthUserSocial{
			UserID: &addUser.ID,
			OpenID: &authQQme.OpenID,
			Source: &qq,
		}
		err = userSocialService.AddUserSocial(userSocial)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		userRole := model.ScaAuthUserRole{
			UserID: uidStr,
			RoleID: enum.User,
		}
		err = userRoleService.AddUserRole(userRole)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		res, data := HandelUserLogin(addUser)
		if !res {
			return
		}
		tokenData, err := json.Marshal(data)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		formattedScript := fmt.Sprintf(script, tokenData, global.CONFIG.System.Web)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedScript))
		return
	} else {
		user, err := userService.QueryUserById(userSocial.UserID)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		res, data := HandelUserLogin(user)
		if !res {
			return
		}
		tokenData, err := json.Marshal(data)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		formattedScript := fmt.Sprintf(script, tokenData, global.CONFIG.System.Web)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(formattedScript))
		return
	}

}

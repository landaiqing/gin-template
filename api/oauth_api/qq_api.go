package oauth_api

import (
	"encoding/json"
	"fmt"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
)

type AuthQQme struct {
	ClientID string `json:"client_id"`
	OpenID   string `json:"openid"`
}

// GetQQRedirectUrl 获取登录地址
// @Summary 获取QQ登录地址
// @Description 获取QQ登录地址
// @Tags 登录
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
		"https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
		clientId, clientSecret, code, redirectURI,
	)
}

// GetQQToken 获取 token
func GetQQToken(url string) (*Token, error) {

	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	// 将响应体解析为 token，并返回
	var token Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// GetQQUserOpenID 获取用户 openid
func GetQQUserOpenID(token *Token) (*AuthQQme, error) {

	// 形成请求
	var userInfoUrl = "https://graph.qq.com/oauth2.0/me" // github用户信息获取接口
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	// 将响应体解析为 AuthQQme，并返回
	var authQQme AuthQQme
	if err = json.NewDecoder(res.Body).Decode(&authQQme); err != nil {
		return nil, err
	}
	return &authQQme, nil
}

// GetQQUserUserInfo 获取用户信息
func GetQQUserUserInfo(token *Token, openId string) (map[string]interface{}, error) {

	clientId := global.CONFIG.OAuth.QQ.ClientID
	// 形成请求
	var userInfoUrl = "https://graph.qq.com/user/get_user_info??access_token=" + token.AccessToken + "&oauth_consumer_key=" + clientId + "&openid=" + openId
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	//req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// QQCallback QQ登录回调
// @Summary QQ登录回调
// @Description QQ登录回调
// @Tags 登录
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
	var token *Token
	if token, err = GetQQToken(tokenAuthUrl); err != nil {
		global.LOG.Error(err)
		return
	}
	authQQme, err := GetQQUserOpenID(token)
	if err != nil {
		return
	}

	// 通过token，获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = GetQQUserUserInfo(token, authQQme.OpenID); err != nil {
		global.LOG.Error(err)
		return
	}
	result.OkWithData(userInfo, c)
	return
}

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
	"time"
)

type GiteeUser struct {
	AvatarURL         string      `json:"avatar_url"`
	Bio               string      `json:"bio"`
	Blog              string      `json:"blog"`
	CreatedAt         time.Time   `json:"created_at"`
	Email             string      `json:"email"`
	EventsURL         string      `json:"events_url"`
	Followers         int         `json:"followers"`
	FollowersURL      string      `json:"followers_url"`
	Following         int         `json:"following"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	HTMLURL           string      `json:"html_url"`
	ID                int         `json:"id"`
	Login             string      `json:"login"`
	Name              string      `json:"name"`
	OrganizationsURL  string      `json:"organizations_url"`
	PublicGists       int         `json:"public_gists"`
	PublicRepos       int         `json:"public_repos"`
	ReceivedEventsURL string      `json:"received_events_url"`
	Remark            string      `json:"remark"`
	ReposURL          string      `json:"repos_url"`
	Stared            int         `json:"stared"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	Type              string      `json:"type"`
	UpdatedAt         time.Time   `json:"updated_at"`
	URL               string      `json:"url"`
	Watched           int         `json:"watched"`
	Weibo             interface{} `json:"weibo"`
}

// GetGiteeRedirectUrl 获取Gitee登录地址
// @Summary 获取Gitee登录地址
// @Description 获取Gitee登录地址
// @Tags OAuth
// @Produce  json
// @Success 200 {string} string "登录地址"
// @Router /api/oauth/gitee/get_url [get]
func (OAuthAPI) GetGiteeRedirectUrl(c *gin.Context) {
	clientID := global.CONFIG.OAuth.Gitee.ClientID
	redirectURI := global.CONFIG.OAuth.Gitee.RedirectURI
	url := "https://gitee.com/oauth/authorize?client_id=" + clientID + "&redirect_uri=" + redirectURI + "&response_type=code"
	result.OkWithData(url, c)
	return
}

// GetGiteeTokenAuthUrl 获取Gitee token
func GetGiteeTokenAuthUrl(code string) string {
	clientId := global.CONFIG.OAuth.Gitee.ClientID
	clientSecret := global.CONFIG.OAuth.Gitee.ClientSecret
	redirectURI := global.CONFIG.OAuth.Gitee.RedirectURI
	return fmt.Sprintf(
		"https://gitee.com/oauth/token?grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s&client_secret=%s",
		code, clientId, redirectURI, clientSecret,
	)
}

// GetGiteeToken 获取 token
func GetGiteeToken(url string) (*Token, error) {

	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodPost, url, nil); err != nil {
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

// GetGiteeUserInfo 获取用户信息
func GetGiteeUserInfo(token *Token) (map[string]interface{}, error) {

	// 形成请求
	var userInfoUrl = "https://gitee.com/api/v5/user" // github用户信息获取接口
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

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GiteeCallback 处理Gitee回调
// @Summary 处理Gitee回调
// @Description 处理Gitee回调
// @Tags OAuth
// @Produce  json
// @Router /api/oauth/gitee/callback [get]
func (OAuthAPI) GiteeCallback(c *gin.Context) {
	var err error
	// 获取 code
	var code = c.Query("code")
	if code == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}
	// 通过 code, 获取 token
	var tokenAuthUrl = GetGiteeTokenAuthUrl(code)
	var token *Token
	if token, err = GetGiteeToken(tokenAuthUrl); err != nil {
		global.LOG.Error(err)
		return
	}

	// 通过token，获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = GetGiteeUserInfo(token); err != nil {
		global.LOG.Error(err)
		return
	}

	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		global.LOG.Error(err)
		return
	}
	var giteeUser GiteeUser
	err = json.Unmarshal(userInfoBytes, &giteeUser)
	if err != nil {
		global.LOG.Error(err)
		return
	}

	Id := strconv.Itoa(giteeUser.ID)
	userSocial, err := userSocialService.QueryUserSocialByUUID(Id, enum.OAuthSourceGitee)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 第一次登录，创建用户
		uid := idgen.NextId()
		uidStr := strconv.FormatInt(uid, 10)
		user := model.ScaAuthUser{
			UID:      &uidStr,
			Username: &giteeUser.Login,
			Nickname: &giteeUser.Name,
			Avatar:   &giteeUser.AvatarURL,
			Blog:     &giteeUser.Blog,
			Email:    &giteeUser.Email,
		}
		addUser, err := userService.AddUser(user)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		gitee := enum.OAuthSourceGitee
		userSocial = model.ScaAuthUserSocial{
			UserID: &addUser.ID,
			UUID:   &Id,
			Source: &gitee,
		}
		err = userSocialService.AddUserSocial(userSocial)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		userRole := model.ScaAuthUserRole{
			UserID: addUser.ID,
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
	}
	return
}

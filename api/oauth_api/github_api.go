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

type GitHubUser struct {
	AvatarURL         string      `json:"avatar_url"`
	Bio               interface{} `json:"bio"`
	Blog              string      `json:"blog"`
	Company           interface{} `json:"company"`
	CreatedAt         string      `json:"created_at"`
	Email             string      `json:"email"`
	EventsURL         string      `json:"events_url"`
	Followers         int         `json:"followers"`
	FollowersURL      string      `json:"followers_url"`
	Following         int         `json:"following"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	GravatarID        string      `json:"gravatar_id"`
	Hireable          interface{} `json:"hireable"`
	HTMLURL           string      `json:"html_url"`
	ID                int         `json:"id"`
	Location          interface{} `json:"location"`
	Login             string      `json:"login"`
	Name              string      `json:"name"`
	NodeID            string      `json:"node_id"`
	NotificationEmail interface{} `json:"notification_email"`
	OrganizationsURL  string      `json:"organizations_url"`
	PublicGists       int         `json:"public_gists"`
	PublicRepos       int         `json:"public_repos"`
	ReceivedEventsURL string      `json:"received_events_url"`
	ReposURL          string      `json:"repos_url"`
	SiteAdmin         bool        `json:"site_admin"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	TwitterUsername   interface{} `json:"twitter_username"`
	Type              string      `json:"type"`
	UpdatedAt         string      `json:"updated_at"`
	URL               string      `json:"url"`
}

// GetRedirectUrl 获取github登录url
// @Summary 获取github登录url
// @Description 获取github登录url
// @Tags Github OAuth
// @Produce  json
// @Success 200 {string} string "登录url"
// @Router /api/oauth/github/get_url [get]
func (OAuthAPI) GetRedirectUrl(c *gin.Context) {
	state := c.Query("state")
	clientId := global.CONFIG.OAuth.Github.ClientID
	redirectUrl := global.CONFIG.OAuth.Github.RedirectURI
	url := "https://github.com/login/oauth/authorize?client_id=" + clientId + "&redirect_uri=" + redirectUrl + "&state=" + state
	result.OkWithData(url, c)
	return
}

// GetTokenAuthUrl 通过code获取token认证url
func GetTokenAuthUrl(code string) string {
	clientId := global.CONFIG.OAuth.Github.ClientID
	clientSecret := global.CONFIG.OAuth.Github.ClientSecret
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		clientId, clientSecret, code,
	)
}

// GetToken 获取 token
func GetToken(url string) (*Token, error) {

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

// GetUserInfo 获取用户信息
func GetUserInfo(token *Token) (map[string]interface{}, error) {

	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口
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

// Callback 登录回调函数
// @Summary 登录回调函数
// @Description 登录回调函数
// @Tags Github OAuth
// @Produce  json
// @Param code query string true "code"
// @Success 200 {string} string "登录成功"
// @Router /api/oauth/github/callback [get]
func (OAuthAPI) Callback(c *gin.Context) {
	var err error
	// 获取 code
	var code = c.Query("code")
	if code == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "ParamsError"), c)
		return
	}

	// 使用channel来接收异步操作的结果
	tokenChan := make(chan *Token)
	userInfoChan := make(chan map[string]interface{})
	errChan := make(chan error)

	// 异步获取token
	go func() {
		var tokenAuthUrl = GetTokenAuthUrl(code)
		token, err := GetToken(tokenAuthUrl)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- token
	}()

	// 异步获取用户信息
	go func() {
		token := <-tokenChan
		if token == nil {
			return
		}
		userInfo, err := GetUserInfo(token)
		if err != nil {
			errChan <- err
			return
		}
		userInfoChan <- userInfo
	}()

	select {
	case err = <-errChan:
		global.LOG.Error(err)
		return
	case userInfo := <-userInfoChan:
		if userInfo == nil {
			global.LOG.Error(<-errChan)
			return
		}
		// 继续处理用户信息
		userInfoBytes, err := json.Marshal(<-userInfoChan)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		var gitHubUser GitHubUser
		err = json.Unmarshal(userInfoBytes, &gitHubUser)
		if err != nil {
			global.LOG.Error(err)
			return
		}
		Id := strconv.Itoa(gitHubUser.ID)
		userSocial, err := userSocialService.QueryUserSocialByUUID(Id, enum.OAuthSourceGithub)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			db := global.DB
			tx := db.Begin() // 开始事务
			if tx.Error != nil {
				global.LOG.Error(tx.Error)
				return
			}
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			// 第一次登录，创建用户
			uid := idgen.NextId()
			uidStr := strconv.FormatInt(uid, 10)
			user := model.ScaAuthUser{
				UID:      &uidStr,
				Username: &gitHubUser.Login,
				Nickname: &gitHubUser.Name,
				Avatar:   &gitHubUser.AvatarURL,
				Blog:     &gitHubUser.Blog,
				Email:    &gitHubUser.Email,
				Gender:   &enum.Male,
			}
			addUser, err := userService.AddUser(user)
			if err != nil {
				tx.Rollback()
				global.LOG.Error(err)
				return
			}
			github := enum.OAuthSourceGithub
			userSocial = model.ScaAuthUserSocial{
				UserID: &uidStr,
				OpenID: &Id,
				Source: &github,
			}
			err = userSocialService.AddUserSocial(userSocial)
			if err != nil {
				tx.Rollback()
				global.LOG.Error(err)
				return
			}
			_, err = global.Casbin.AddRoleForUser(uidStr, enum.User)
			if err != nil {
				tx.Rollback()
				global.LOG.Error(err)
				return
			}
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				global.LOG.Error(err)
				return
			}
			HandleLoginResponse(c, *addUser.UID)
		} else {
			HandleLoginResponse(c, *userSocial.UserID)
		}
		return
	}
}

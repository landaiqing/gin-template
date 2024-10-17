package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
	"time"
)

// ResponseData 返回数据
type ResponseData struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	UID          *string  `json:"uid"`
	UserInfo     UserInfo `json:"user_info"`
}
type UserInfo struct {
	Username string    `json:"username,omitempty"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	Phone    string    `json:"phone,omitempty"`
	Email    string    `json:"email,omitempty"`
	Gender   string    `json:"gender"`
	Status   int64     `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

// SetSession sets session data with key and data
func SetSession(c *gin.Context, key string, data interface{}) error {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("SetSession failed: ", err)
		return err
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		global.LOG.Error("SetSession failed: ", err)
		return err
	}
	session.Values[key] = jsonData
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		global.LOG.Error("SetSession failed: ", err)
		return err
	}
	return nil
}

// GetSession gets session data with key
func GetSession(c *gin.Context, key string) *ResponseData {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return nil
	}
	jsonData, ok := session.Values[key]
	if !ok {
		global.LOG.Error("GetSession failed: ", "key not found")
		return nil
	}
	data := ResponseData{}
	err = json.Unmarshal(jsonData.([]byte), &data)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return nil
	}
	return &data
}

// DelSession deletes session data with key
func DelSession(c *gin.Context, key string) {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("DelSession failed: ", err)
		return
	}
	delete(session.Values, key)
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		global.LOG.Error("DelSession failed: ", err)
		return
	}
	return
}

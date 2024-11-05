package utils

import (
	"encoding/gob"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"schisandra-cloud-album/global"
)

// SessionData 返回数据
type SessionData struct {
	RefreshToken string `json:"refresh_token"`
	UID          string `json:"uid"`
}

// SetSession sets session data with key and data
func SetSession(c *gin.Context, key string, data SessionData) error {
	gob.Register(SessionData{})
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
func GetSession(c *gin.Context, key string) SessionData {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return SessionData{}
	}
	jsonData, ok := session.Values[key]
	if !ok {
		global.LOG.Error("GetSession failed: ", "key not found")
		return SessionData{}
	}
	data := SessionData{}
	err = json.Unmarshal(jsonData.([]byte), &data)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return SessionData{}
	}
	return data
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

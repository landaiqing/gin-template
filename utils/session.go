package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/api/user_api/dto"
	"schisandra-cloud-album/global"
)

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
func GetSession(c *gin.Context, key string) dto.ResponseData {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return dto.ResponseData{}
	}
	jsonData, ok := session.Values[key]
	if !ok {
		global.LOG.Error("GetSession failed: ", "key not found")
		return dto.ResponseData{}
	}
	var data dto.ResponseData
	err = json.Unmarshal(jsonData.([]byte), &data)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return dto.ResponseData{}
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

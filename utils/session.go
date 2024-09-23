package utils

import (
	"encoding/gob"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"schisandra-cloud-album/global"
)

func SetSession(c *gin.Context, key string, data interface{}) error {
	session, err := global.Session.Get(c.Request, key)
	if err != nil {
		global.LOG.Error("SetSession failed: ", err)
		return err
	}
	gob.Register(data)
	session.Values[key] = data
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		global.LOG.Error("SetSession failed: ", err)
		return err
	}
	return nil
}

func GetSession(c *gin.Context, key string) interface{} {
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
	var data interface{}
	err = json.Unmarshal(jsonData.([]byte), &data)
	if err != nil {
		global.LOG.Error("GetSession failed: ", err)
		return nil
	}
	return data
}

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

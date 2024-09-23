package utils

import (
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

package core

import (
	"github.com/sirupsen/logrus"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/middleware"
	"time"
)

func InitLogger() {
	// 按照日期格式化输出
	mylog := middleware.NewDateLog(&middleware.DateLogConfig{
		Date: time.Now().Format("2006-01-02"),
		Path: global.CONFIG.Logger.Director,
		Name: global.CONFIG.Logger.LogName,
	})
	log := mylog.Init()
	level, err := logrus.ParseLevel(global.CONFIG.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level) //设置日志级别
	global.LOG = log
}

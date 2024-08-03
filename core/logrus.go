package core

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"schisandra-cloud-album/global"
)

const (
	red    = 31
	yellow = 33
	blue   = 34
	gray   = 37
)

type LogFormatter struct{}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level显示不同的颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	log := global.CONFIG.Logger
	// 自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s %s : %s\n", log.Prefix, timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m : %s\n", log.Prefix, timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	newLog := logrus.New()
	newLog.SetOutput(os.Stdout)                           //设置输出类型
	newLog.SetReportCaller(global.CONFIG.Logger.ShowLine) //设置是否显示函数名和行号
	newLog.SetFormatter(&LogFormatter{})                  //设置日志格式
	level, err := logrus.ParseLevel(global.CONFIG.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	newLog.SetLevel(level) //设置日志级别
	global.LOG = newLog
	InitDefaultLogger()
	return newLog
}

func InitDefaultLogger() {
	//全局日志
	logrus.SetOutput(os.Stdout)                           //设置输出类型
	logrus.SetReportCaller(global.CONFIG.Logger.ShowLine) //设置是否显示函数名和行号
	logrus.SetFormatter(&LogFormatter{})                  //设置日志格式
	level, err := logrus.ParseLevel(global.CONFIG.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level) //设置日志级别
}

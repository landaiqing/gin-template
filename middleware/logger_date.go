package middleware

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"schisandra-cloud-album/global"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ccRed    = 1
	ccYellow = 3
	ccBlue   = 4
	ccCyan   = 6
	ccGray   = 7
)

type LogConsoleFormatter struct{}

func (s *LogConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var color int
	switch entry.Level {

	case logrus.ErrorLevel:
		color = ccRed
	case logrus.WarnLevel:
		color = ccYellow
	case logrus.InfoLevel:
		color = ccBlue
	case logrus.DebugLevel:
		color = ccCyan
	default:
		color = ccGray
	}
	// 设置 buffer 缓冲区
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	fileVal := fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
	fmt.Fprintf(b, "[%s] \033[3%dm[%s]\033[0m [%s] %s %s\n", global.CONFIG.Logger.Prefix, color, entry.Level, timestamp, fileVal, entry.Message)
	return b.Bytes(), nil
}

type LogFileFormatter struct{}

func (s *LogFileFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var length int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		length = entry.Caller.Line
	}
	msg := fmt.Sprintf("[%s] %s [%s:%d] %s\n", strings.ToUpper(entry.Level.String()), timestamp, file, length, entry.Message)
	return []byte(msg), nil
}

type ConsoleHook struct {
	formatter logrus.Formatter
}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	originalFormatter := entry.Logger.Formatter
	entry.Logger.Formatter = hook.formatter
	defer func() { entry.Logger.Formatter = originalFormatter }()
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(line)
	return err
}

type FileHook struct {
	formatter logrus.Formatter
	file      *os.File
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	originalFormatter := entry.Logger.Formatter
	entry.Logger.Formatter = hook.formatter
	defer func() { entry.Logger.Formatter = originalFormatter }()
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.file.Write(line)
	return err
}

type DateLogConfig struct {
	Date string
	Path string
	Name string
}

func NewDateLog(d *DateLogConfig) *DateLogConfig {
	return &DateLogConfig{
		Date: d.Date,
		Path: d.Path,
		Name: d.Name,
	}
}

func (d *DateLogConfig) Init() *logrus.Logger {
	// 实例化 logrus
	log := logrus.New()
	// 设置是否输出文件名和行号信息
	log.SetReportCaller(global.CONFIG.Logger.ShowLine)
	// 将 logrus 的默认输出丢弃，确保日志只通过 hooks 输出
	log.SetOutput(io.Discard)

	// 控制台输出的 hook
	consoleHook := &ConsoleHook{
		formatter: &LogConsoleFormatter{},
	}
	// 添加控制台输出的 hook
	log.AddHook(consoleHook)

	// 文件路径
	filename := fmt.Sprintf("%s/%s/%s.log", d.Path, d.Date, d.Name)
	// 创建目录
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", d.Path, d.Date), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	// 打开文件，如果文件不存在，则创建文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}

	// 文件输出的 hook
	fileHook := &FileHook{
		formatter: &LogFileFormatter{},
		file:      file,
	}
	// 添加文件输出的 hook
	log.AddHook(fileHook)
	return log
}

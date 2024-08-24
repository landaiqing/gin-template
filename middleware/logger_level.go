package middleware

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	allLog   = "all"
	errLog   = "err"
	warnLog  = "warn"
	infoLog  = "info"
	debugLog = "debug"
)

type LevelFormatter struct{}

func (l *LevelFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	msg := fmt.Sprintf("[%s] %s [%s:%d] %s\n", strings.ToUpper(entry.Level.String()), timestamp, file, len, entry.Message)
	return []byte(msg), nil
}

type LevelConfig struct {
	Date string
	Name string
	Path string
}

func NewLevelLog(d *LevelConfig) *LevelConfig {
	return &LevelConfig{
		Date: d.Date,
		Path: d.Path,
		Name: d.Name,
	}
}

func (l *LevelConfig) Init() *logrus.Logger {
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetOutput(io.Discard)

	err := os.MkdirAll(fmt.Sprintf("%s/%s", l.Path, l.Date), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	allFile, _ := os.OpenFile(fmt.Sprintf("%s/%s/%s-%s.log", l.Path, l.Date, l.Name, allLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	errFile, _ := os.OpenFile(fmt.Sprintf("%s/%s/%s-%s.log", l.Path, l.Date, l.Name, errLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	warnFile, _ := os.OpenFile(fmt.Sprintf("%s/%s/%s-%s.log", l.Path, l.Date, l.Name, warnLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	infoFile, _ := os.OpenFile(fmt.Sprintf("%s/%s/%s-%s.log", l.Path, l.Date, l.Name, infoLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	debugFile, _ := os.OpenFile(fmt.Sprintf("%s/%s/%s-%s.log", l.Path, l.Date, l.Name, debugLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)

	fileHook := &LevelHook{
		formatter: &LevelFormatter{},
		file:      allFile,
		errFile:   errFile,
		warnFile:  warnFile,
		infoFile:  infoFile,
		debugFile: debugFile,
	}

	log.AddHook(fileHook)
	return log
}

type LevelHook struct {
	formatter logrus.Formatter
	file      *os.File
	errFile   *os.File
	warnFile  *os.File
	infoFile  *os.File
	debugFile *os.File
}

func (l *LevelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *LevelHook) Fire(entry *logrus.Entry) error {
	originalFormatter := entry.Logger.Formatter
	entry.Logger.Formatter = l.formatter
	defer func() { entry.Logger.Formatter = originalFormatter }()
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = l.file.Write(line)

	switch entry.Level {
	case logrus.ErrorLevel:
		_, err = l.errFile.Write(line)
	case logrus.WarnLevel:
		_, err = l.warnFile.Write(line)
	case logrus.InfoLevel:
		_, err = l.infoFile.Write(line)
	case logrus.DebugLevel:
		_, err = l.debugFile.Write(line)
	default:
		panic("unhandled default case")
	}
	return err
}

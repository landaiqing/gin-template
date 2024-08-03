package config

// Logger 配置
type Logger struct {
	Level        string `yaml:"level"`          // 日志级别
	Prefix       string `yaml:"prefix"`         // 日志前缀
	Director     string `yaml:"director"`       // 日志文件存放目录
	ShowLine     bool   `yaml:"show-line"`      // 是否显示文件行号
	LogInConsole string `yaml:"log-in-console"` // 是否在控制台打印日志
}

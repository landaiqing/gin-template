package config

// Logger 配置
type Logger struct {
	Level        string `yaml:"level"`        // 日志级别
	Prefix       string `yaml:"prefix"`       // 日志前缀
	Director     string `yaml:"director"`     // 日志文件存放目录
	ShowLine     string `yaml:"showLine"`     // 是否显示文件行号
	LogInConsole string `yaml:"logInConsole"` // 是否在控制台打印日志
}

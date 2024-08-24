package config

type Casbin struct {
	// 权限模型文件路径
	ModelPath string `yaml:"model-path"`
	// 数据库前缀
	TablePrefix string `yaml:"table-prefix"`
	// 数据库表明
	TableName string `yaml:"table-name"`
}

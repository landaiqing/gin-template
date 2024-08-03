package config

// System 系统配置
type System struct {
	Host string `yaml:"host"` //主机地址
	Port string `yaml:"port"` //端口号
	Env  string `yaml:"env"`  //环境
}

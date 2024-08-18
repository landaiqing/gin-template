package config

import "fmt"

// System 系统配置
type System struct {
	Host string `yaml:"host"` //主机地址
	Port string `yaml:"port"` //端口号
	Env  string `yaml:"env"`  //环境
	Web  string `yaml:"web"`  //web地址
}

func (s *System) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

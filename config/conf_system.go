package config

import "fmt"

// System 系统配置
type System struct {
	Host     string `yaml:"host"`     //主机地址
	Port     string `yaml:"port"`     //端口号
	Env      string `yaml:"env"`      //环境
	Protocol string `yaml:"protocol"` //协议
	Web      string `yaml:"web"`      //web地址
	Ip       string `yaml:"ip"`       //ip地址
}

func (s *System) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *System) WebURL() string {
	return fmt.Sprintf("%s://%s", s.Protocol, s.Web)
}

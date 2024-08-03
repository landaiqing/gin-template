package config

import "strconv"

// MySQL 配置
type MySQL struct {
	Host          string `yaml:"host"`            //主机地址
	Port          int    `yaml:"port"`            //端口号
	DB            string `yaml:"db"`              //数据库名
	User          string `yaml:"user"`            //用户名
	Password      string `yaml:"password"`        //密码
	LogLevel      string `yaml:"log-level"`       //日志级别 debug: 输出全部SQL语句; release: 只输出错误信息
	Config        string `yaml:"config"`          //高级配置
	MaxIdleConnes int    `yaml:"max-idle-connes"` //最大空闲连接数
	MaxOpenConnes int    `yaml:"max-open-connes"` //最大连接数
}

func (m *MySQL) Dsn() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DB + "?" + m.Config
}

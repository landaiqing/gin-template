package config

// MySQL 配置
type MySQL struct {
	Host     string `yaml:"host"`      //主机地址
	Port     string `yaml:"port"`      //端口号
	DB       string `yaml:"db"`        //数据库名
	User     string `yaml:"user"`      //用户名
	Password string `yaml:"password"`  //密码
	LogLevel string `yaml:"log-level"` // 日志级别 debug: 输出全部SQL语句; release: 只输出错误信息
}

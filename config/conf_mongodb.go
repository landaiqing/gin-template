package config

import "strconv"

type MongoDB struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	AuthSource  string `yaml:"auth-source"`
	DB          string `yaml:"db"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	MaxOpenConn uint64 `yaml:"max-open-conn"`
	MaxIdleConn uint64 `yaml:"max-idle-conn"`
	Timeout     int    `yaml:"timeout"`
}

func (m *MongoDB) MongoDsn() string {
	return "mongodb://" + m.Host + ":" + strconv.Itoa(m.Port)
}

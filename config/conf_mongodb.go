package config

import "strconv"

type MongoDB struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	AuthSource  string `yaml:"auth-source"`
	DB          string `yaml:"db"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	MaxOpenConn int    `yaml:"max-open-conn"`
	MaxIdleConn int    `yaml:"max-idle-conn"`
}

func (m *MongoDB) MongoDsn() string {
	return "mongodb://" + m.Host + ":" + strconv.Itoa(m.Port)
}

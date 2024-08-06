package config

import "time"

type Redis struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Password    string        `yaml:"password"`
	Db          int           `yaml:"db"`
	MaxActive   int           `yaml:"max-active"`
	MaxIdle     int           `yaml:"max-idle"`
	MinIdle     int           `yaml:"min-idle"`
	PoolSize    int           `yaml:"pool-size"`
	PoolTimeout time.Duration `yaml:"pool-timeout"`
}

func (r *Redis) Addr() string {
	return r.Host + ":" + r.Port
}

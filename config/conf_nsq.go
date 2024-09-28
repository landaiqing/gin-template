package config

import "fmt"

type NSQ struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (n *NSQ) Addr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

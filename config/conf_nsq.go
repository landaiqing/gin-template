package config

import "fmt"

type NSQ struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	LookupdHost string `yaml:"lookupd-host"`
	LookupdPort int    `yaml:"lookupd-port"`
}

func (n *NSQ) NsqAddr() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (n *NSQ) LookupdAddr() string {
	return fmt.Sprintf("%s:%d", n.LookupdHost, n.LookupdPort)
}

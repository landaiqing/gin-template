package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"schisandra-cloud-album/config"
	"schisandra-cloud-album/global"
)

// InitConfig 初始化配置
func InitConfig() {
	const ConfigFile = "config.yaml"
	c := &config.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("get yaml config error: %s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatal("config init unmarshal error: ", err)
	}
	global.CONFIG = c
}

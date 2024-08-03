package main

import (
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
)

func main() {
	// 初始化配置
	core.InitConfig()
	core.InitLogger()
	core.InitGorm()
	global.LOG.Error("hello world")
}

package main

import (
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router"
)

func main() {
	// 初始化配置
	core.InitConfig()
	core.InitLogger()
	core.InitGorm()
	r := router.InitRouter()
	addr := global.CONFIG.System.Addr()
	global.LOG.Info("Server run on ", addr)
	err := r.Run(addr)
	if err != nil {
		return
	}
}

package main

import (
	"schisandra-cloud-album/cmd"
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router"
)

func main() {
	// 初始化配置
	core.InitConfig()
	core.InitLogger()
	core.InitGorm()
	// 命令行参数绑定
	option := cmd.Parse()
	if cmd.IsStopWeb(&option) {
		cmd.SwitchOption(&option)
		return
	}
	r := router.InitRouter()
	addr := global.CONFIG.System.Addr()
	global.LOG.Info("Server run on ", addr)
	err := r.Run(addr)
	if err != nil {
		global.LOG.Fatalf(err.Error())
	}
}

package main

import (
	"schisandra-cloud-album/cmd"
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router"
)

func main() {
	// 初始化配置
	core.InitConfig()  // 读取配置文件
	core.InitLogger()  // 初始化日志
	core.InitGorm()    // 初始化数据库
	core.InitRedis()   // 初始化redis
	core.InitCaptcha() // 初始化验证码
	// 命令行参数绑定
	option := cmd.Parse()
	if cmd.IsStopWeb(&option) {
		cmd.SwitchOption(&option)
		return
	}
	r := router.InitRouter() // 初始化路由
	addr := global.CONFIG.System.Addr()
	global.LOG.Info("Server run on ", addr)
	err := r.Run(addr)
	if err != nil {
		global.LOG.Fatalf(err.Error())
	}
}

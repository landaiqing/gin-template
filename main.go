package main

import (
	"schisandra-cloud-album/cmd"
	"schisandra-cloud-album/core"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/router"
)

func main() {

	// 初始化配置
	core.InitConfig()      // 读取配置文件
	core.InitLogger()      // 初始化日志
	core.InitGorm()        // 初始化数据库
	core.InitRedis()       // 初始化redis
	core.InitCaptcha()     // 初始化验证码
	core.InitIDGenerator() // 初始化ID生成器
	core.InitWechat()      // 初始化微信
	core.InitCasbin()      // 初始化Casbin
	core.InitIP2Region()   // 初始化IP2Region
	// 命令行参数绑定
	option := cmd.Parse()
	if cmd.IsStopWeb(&option) {
		cmd.SwitchOption(&option)
		return
	}
	r := router.InitRouter() // 初始化路由
	addr := global.CONFIG.System.Addr()

	err := r.Run(addr)
	if err != nil {
		global.LOG.Fatalf(err.Error())
	}
}

package cmd

import (
	"flag"
	"schisandra-cloud-album/cmd/db"
	"schisandra-cloud-album/global"
)

type Option struct {
	DB bool
}

// Parse 解析命令行： go run main.go -db
func Parse() Option {
	// go run main.go -db
	DB := flag.Bool("db", false, "初始化数据库")
	// 解析命令
	flag.Parse()
	return Option{
		DB: *DB,
	}
}

func IsStopWeb(option *Option) bool {
	if option.DB {
		global.LOG.Infof("停止web项目")
		return true
	}
	return false // 停止web项目
}

func SwitchOption(option *Option) {
	if option.DB {
		// 迁移数据库
		db.MakeMigration()
		return
	}
	flag.Usage()
}

package core

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"schisandra-cloud-album/global"
	"time"
)

func InitGorm() {
	global.DB = mySQlConnect()
}
func mySQlConnect() *gorm.DB {
	if global.CONFIG.MySQL.Host == "" {
		return nil
	}
	dsn := global.CONFIG.MySQL.Dsn()
	var mysqlLogger logger.Interface
	if global.CONFIG.System.Env == "dev" {
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(global.CONFIG.MySQL.MaxIdleConnes)
	sqlDB.SetMaxOpenConns(global.CONFIG.MySQL.MaxOpenConnes)
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //连接最大复用时间
	return db
}

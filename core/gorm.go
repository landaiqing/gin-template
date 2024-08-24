package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"schisandra-cloud-album/global"
	"time"
)

func InitGorm() {
	global.DB = MySQlConnect()
}
func MySQlConnect() *gorm.DB {
	if global.CONFIG.MySQL.Host == "" {
		global.LOG.Warnln("未配置MySQL,取消连接！")
		return nil
	}
	dsn := global.CONFIG.MySQL.Dsn()
	var mysqlLogger logger.Interface
	if global.CONFIG.System.Env == "debug" {
		//自定义日子模板 打印SQL语句
		mysqlLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, //慢sql日志
				LogLevel:                  logger.Info, //级别
				Colorful:                  true,        //颜色
				IgnoreRecordNotFoundError: true,        //忽略RecordNotFoundError
				ParameterizedQueries:      true,        //格式化SQL语句

			})
	} else {
		logfile, _ := os.OpenFile("/tmp/logs/mysql.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		mysqlLogger = logger.New(
			log.New(logfile, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,  //慢sql日志
				LogLevel:                  logger.Error, //级别
				Colorful:                  true,         //颜色
				IgnoreRecordNotFoundError: true,         //忽略RecordNotFoundError
				ParameterizedQueries:      true,         //格式化SQL语句
			})
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      mysqlLogger,
		PrepareStmt: true,
	})
	if err != nil {
		global.LOG.Fatalf(fmt.Sprintf("[%s] MySQL 连接失败", dsn))
	}
	sqlDB, err := db.DB()
	if err != nil {
		global.LOG.Fatalf(fmt.Sprintf("[%s] MySQL 获取DB失败", dsn))
	}
	sqlDB.SetMaxIdleConns(global.CONFIG.MySQL.MaxIdleConnes)
	sqlDB.SetMaxOpenConns(global.CONFIG.MySQL.MaxOpenConnes)
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //连接最大复用时间
	return db
}

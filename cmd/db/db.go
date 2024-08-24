package db

import (
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/model"
)

func MakeMigration() {
	var err error
	global.LOG.Infof("开始迁移数据库")
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.ScaAuthUser{},
		&model.ScaAuthPermission{},
		&model.ScaAuthRole{},
		&model.ScaAuthUserDevice{},
		&model.ScaAuthUserSocial{},
	)
	if err != nil {
		global.LOG.Error("数据库迁移失败: %v", err)
		return
	}
	global.LOG.Info("数据库迁移成功")

}

package db

import "schisandra-cloud-album/global"

func MakeMigration() {
	var err error
	global.LOG.Infof("开始迁移数据库")
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
	//&models.ImageModel{},
	//&models.TagModel{},
	//&models.MessageModel{},
	//&models.AdvertModel{},
	//&models.UserModel{},
	//&models.CommentModel{},
	//&models.ArticleModel{},
	//&models.MenuModel{},
	//&models.MenuImageModel{},
	//&models.FeedbackModel{},
	//&models.LoginDataModel{},
	)
	if err != nil {
		global.LOG.Error("数据库迁移失败: %v", err)
		return
	}
	global.LOG.Info("数据库迁移成功")

}

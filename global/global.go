package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"schisandra-cloud-album/config"
)

// Config 全局配置文件
var (
	CONFIG *config.Config
	DB     *gorm.DB
	LOG    *logrus.Logger
)

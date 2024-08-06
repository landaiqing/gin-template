package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"gorm.io/gorm"
	"schisandra-cloud-album/config"
)

// Config 全局配置文件
var (
	CONFIG        *config.Config
	DB            *gorm.DB
	LOG           *logrus.Logger
	TextCaptcha   click.Captcha
	SlideCaptcha  slide.Captcha
	RotateCaptcha rotate.Captcha
	REDIS         *redis.Client
)

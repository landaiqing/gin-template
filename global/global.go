package global

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/casbin/casbin/v2"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/nsqio/go-nsq"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"github.com/zmexing/go-sensitive-word"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"schisandra-cloud-album/config"
)

// Config 全局配置文件
var (
	CONFIG             *config.Config                   // 配置文件
	DB                 *gorm.DB                         // 数据库连接
	LOG                *logrus.Logger                   // 日志
	TextCaptcha        click.Captcha                    // 文本验证码
	LightTextCaptcha   click.Captcha                    // 亮色文本验证码
	ClickShapeCaptcha  click.Captcha                    // 点击形状验证码
	SlideCaptcha       slide.Captcha                    // 滑块验证码
	RotateCaptcha      rotate.Captcha                   // 旋转验证码
	SlideRegionCaptcha slide.Captcha                    // 滑块区域验证码
	REDIS              *redis.Client                    // redis连接
	Wechat             *officialAccount.OfficialAccount // 微信公众号
	Casbin             *casbin.CachedEnforcer           // casbin权限管理器
	IP2Location        *xdb.Searcher                    // IP地址定位
	MongoDB            *mongo.Client                    // MongoDB连接
	Session            *redisstore.RedisStore           // session存储
	NSQProducer        *nsq.Producer                    // NSQ生产者
	SensitiveManager   *go_sensitive_word.Manager       // 敏感词管理器
)

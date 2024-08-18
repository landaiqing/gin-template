package core

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"os"
	"schisandra-cloud-album/global"
)

func InitWechat() {
	OfficialAccountApp, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  global.CONFIG.Wechat.AppID,
		Secret: global.CONFIG.Wechat.AppSecret,
		Token:  global.CONFIG.Wechat.Token,
		AESKey: global.CONFIG.Wechat.AESKey,
		//Log: officialAccount.Log{
		//	Level: "debug",
		//	File:  "./wechat.log",
		//},
		ResponseType: os.Getenv("response_type"),
		HttpDebug:    true,
		Debug:        true,
		Cache: kernel.NewRedisClient(&kernel.UniversalOptions{
			Addrs:    []string{global.CONFIG.Redis.Addr()},
			Password: global.CONFIG.Redis.Password,
			DB:       2,
		}),
	})
	if err != nil {
		panic(err)
	}
	global.Wechat = OfficialAccountApp
}

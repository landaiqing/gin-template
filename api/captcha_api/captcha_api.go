package captcha_api

import (
	"encoding/json"
	"fmt"
	"github.com/wenlng/go-captcha/v2/base/option"
	"log"
	"schisandra-cloud-album/global"
)

// GenerateTextCaptcha 生成文本验证码
func GenerateTextCaptcha() {
	captData, err := global.TextCaptcha.Generate()
	if err != nil {
		log.Fatalln(err)
	}

	dotData := captData.GetData()
	if dotData == nil {
		log.Fatalln(">>>>> generate err")
	}

	dots, _ := json.Marshal(dotData)
	fmt.Println(">>>>> ", string(dots))

	err = captData.GetMasterImage().SaveToFile("./.caches/master.jpg", option.QualityNone)
	if err != nil {
		fmt.Println(err)
	}
	err = captData.GetThumbImage().SaveToFile("./.caches/thumb.png")
	if err != nil {
		fmt.Println(err)
	}
}

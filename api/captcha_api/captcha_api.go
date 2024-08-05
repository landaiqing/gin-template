package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
)

var textCapt click.Captcha

func init() {

}

func main() {
	captData, err := textCapt.Generate()
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

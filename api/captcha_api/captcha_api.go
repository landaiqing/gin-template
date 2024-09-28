package captcha_api

import (
	"encoding/json"
	"fmt"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha-assets/helper"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"log"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"strconv"
	"strings"
	"time"
)

// GenerateRotateCaptcha 生成旋转验证码
// @Summary 生成旋转验证码
// @Description 生成旋转验证码
// @Tags 旋转验证码
// @Success 200 {string} json
// @Router /api/captcha/rotate/get [get]
func (CaptchaAPI) GenerateRotateCaptcha(c *gin.Context) {
	captchaData, err := global.RotateCaptcha.Generate()
	if err != nil {
		global.LOG.Fatalln(err)
		result.FailWithNull(c)
		return
	}
	blockData := captchaData.GetData()
	if blockData == nil {
		result.FailWithNull(c)
		return
	}

	masterImageBase64 := captchaData.GetMasterImage().ToBase64()
	thumbImageBase64 := captchaData.GetThumbImage().ToBase64()
	dotsByte, err := json.Marshal(blockData)
	if err != nil {
		global.LOG.Fatalln(err)
		result.FailWithNull(c)
		return
	}

	key := helper.StringToMD5(string(dotsByte))
	err = redis.Set(constant.UserLoginCaptchaRedisKey+key, dotsByte, time.Minute).Err()
	if err != nil {
		global.LOG.Fatalln(err)
		result.FailWithNull(c)
		return
	}

	result.OkWithData(map[string]interface{}{
		"key":   key,
		"image": masterImageBase64,
		"thumb": thumbImageBase64,
	}, c)
}

// CheckRotateData 验证旋转验证码
// @Summary 验证旋转验证码
// @Description 验证旋转验证码
// @Tags 旋转验证码
// @Param angle query string true "验证码角度"
// @Param key query string true "验证码key"
// @Success 200 {string} json
// @Router /api/captcha/rotate/check [post]
func (CaptchaAPI) CheckRotateData(c *gin.Context) {
	var rotateRequest RotateCaptchaRequest
	if err := c.ShouldBindJSON(&rotateRequest); err != nil {
		result.FailWithNull(c)
		return
	}

	cacheDataByte, err := redis.Get(constant.UserLoginCaptchaRedisKey + rotateRequest.Key).Bytes()
	if err != nil || len(cacheDataByte) == 0 {
		result.FailWithCodeAndMessage(1011, ginI18n.MustGetMessage(c, "CaptchaExpired"), c)
		return
	}

	var dct *rotate.Block
	if err := json.Unmarshal(cacheDataByte, &dct); err != nil {
		result.FailWithNull(c)
		return
	}

	sAngle, err := strconv.ParseFloat(fmt.Sprintf("%v", rotateRequest.Angle), 64)
	if err != nil {
		result.FailWithNull(c)
		return
	}

	chkRet := rotate.CheckAngle(int64(sAngle), int64(dct.Angle), 2)
	if chkRet {
		result.OkWithMessage("success", c)
		return
	}

	result.FailWithMessage("fail", c)
}

// GenerateBasicTextCaptcha 生成基础文字验证码
// @Summary 生成基础文字验证码
// @Description 生成基础文字验证码
// @Tags 基础文字验证码
// @Param type query string true "验证码类型"
// @Success 200 {string} json
// @Router /api/captcha/text/get [get]
func (CaptchaAPI) GenerateBasicTextCaptcha(c *gin.Context) {
	var capt click.Captcha
	if c.Query("type") == "light" {
		capt = global.LightTextCaptcha
	} else {
		capt = global.TextCaptcha
	}
	captData, err := capt.Generate()
	if err != nil {
		global.LOG.Fatalln(err)
	}
	dotData := captData.GetData()
	if dotData == nil {
		result.FailWithNull(c)
		return
	}
	var masterImageBase64, thumbImageBase64 string
	masterImageBase64 = captData.GetMasterImage().ToBase64()
	thumbImageBase64 = captData.GetThumbImage().ToBase64()

	dotsByte, err := json.Marshal(dotData)
	if err != nil {
		result.FailWithNull(c)
		return
	}
	key := helper.StringToMD5(string(dotsByte))
	err = redis.Set("user:login:client:"+key, dotsByte, time.Minute).Err()
	if err != nil {
		result.FailWithNull(c)
		return
	}
	bt := map[string]interface{}{
		"key":   key,
		"image": masterImageBase64,
		"thumb": thumbImageBase64,
	}
	result.OkWithData(bt, c)
}

// CheckClickData 验证基础文字验证码
// @Summary 验证基础文字验证码
// @Description 验证基础文字验证码
// @Tags 基础文字验证码
// @Param captcha query string true "验证码"
// @Param key query string true "验证码key"
// @Success 200 {string} json
// @Router /api/captcha/text/check [get]
func (CaptchaAPI) CheckClickData(c *gin.Context) {
	dots := c.Query("dots")
	key := c.Query("key")
	if dots == "" || key == "" {
		result.FailWithNull(c)
		return
	}
	cacheDataByte, err := redis.Get("user:login:client:" + key).Bytes()
	if len(cacheDataByte) == 0 || err != nil {
		result.FailWithNull(c)
		return
	}
	src := strings.Split(dots, ",")

	var dct map[int]*click.Dot
	if err := json.Unmarshal(cacheDataByte, &dct); err != nil {
		result.FailWithNull(c)
		return
	}
	chkRet := false
	if (len(dct) * 2) == len(src) {
		for i := 0; i < len(dct); i++ {
			dot := dct[i]
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)

			chkRet = click.CheckPoint(int64(sx), int64(sy), int64(dot.X), int64(dot.Y), int64(dot.Width), int64(dot.Height), 0)
			if !chkRet {
				break
			}
		}
	}
	if chkRet {
		result.OkWithMessage("success", c)
		return
	}
	result.FailWithMessage("fail", c)
}

// GenerateClickShapeCaptcha 生成点击形状验证码
// @Summary 生成点击形状验证码
// @Description 生成点击形状验证码
// @Tags 点击形状验证码
// @Success 200 {string} json
// @Router /api/captcha/shape/get [get]
func (CaptchaAPI) GenerateClickShapeCaptcha(c *gin.Context) {
	captData, err := global.ClickShapeCaptcha.Generate()
	if err != nil {
		log.Fatalln(err)
	}
	dotData := captData.GetData()
	if dotData == nil {
		result.FailWithNull(c)
		return
	}
	var masterImageBase64, thumbImageBase64 string
	masterImageBase64 = captData.GetMasterImage().ToBase64()
	thumbImageBase64 = captData.GetThumbImage().ToBase64()

	dotsByte, err := json.Marshal(dotData)
	if err != nil {
		result.FailWithNull(c)
		return
	}
	key := helper.StringToMD5(string(dotsByte))
	err = redis.Set(key, dotsByte, time.Minute).Err()
	if err != nil {
		result.FailWithNull(c)
		return
	}
	bt := map[string]interface{}{
		"key":   key,
		"image": masterImageBase64,
		"thumb": thumbImageBase64,
	}
	result.OkWithData(bt, c)
}

// GenerateSlideBasicCaptData 滑块基础验证码
// @Summary 滑块基础验证码
// @Description 滑块基础验证码
// @Tags 滑块基础验证码
// @Success 200 {string} json
// @Router /api/captcha/slide/generate [get]
func (CaptchaAPI) GenerateSlideBasicCaptData(c *gin.Context) {
	captData, err := global.SlideCaptcha.Generate()
	if err != nil {
		global.LOG.Fatalln(err)
	}
	blockData := captData.GetData()
	if blockData == nil {
		result.FailWithNull(c)
		return
	}
	var masterImageBase64, tileImageBase64 string
	masterImageBase64 = captData.GetMasterImage().ToBase64()

	tileImageBase64 = captData.GetTileImage().ToBase64()

	dotsByte, err := json.Marshal(blockData)
	if err != nil {
		result.FailWithNull(c)
		return
	}
	key := helper.StringToMD5(string(dotsByte))
	err = redis.Set(constant.CommentSubmitCaptchaRedisKey+key, dotsByte, time.Minute).Err()
	if err != nil {
		result.FailWithNull(c)
		return
	}
	bt := map[string]interface{}{
		"key":          key,
		"image":        masterImageBase64,
		"thumb":        tileImageBase64,
		"thumb_width":  blockData.Width,
		"thumb_height": blockData.Height,
		"thumb_x":      blockData.TileX,
		"thumb_y":      blockData.TileY,
	}
	result.OkWithData(bt, c)
}

// GenerateSlideRegionCaptData 生成滑动区域形状验证码
// @Summary 生成滑动区域形状验证码
// @Description 生成滑动区域形状验证码
// @Tags 生成滑动区域形状验证码
// @Success 200 {string} json
// @Router /api/captcha/shape/slide/region/get [get]
func (CaptchaAPI) GenerateSlideRegionCaptData(c *gin.Context) {
	captData, err := global.SlideRegionCaptcha.Generate()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	blockData := captData.GetData()
	if blockData == nil {
		result.FailWithNull(c)
		return
	}

	var masterImageBase64, tileImageBase64 string
	masterImageBase64 = captData.GetMasterImage().ToBase64()
	tileImageBase64 = captData.GetTileImage().ToBase64()

	blockByte, err := json.Marshal(blockData)
	if err != nil {
		result.FailWithNull(c)
		return
	}
	key := helper.StringToMD5(string(blockByte))
	err = redis.Set(key, blockByte, time.Minute).Err()
	if err != nil {
		result.FailWithNull(c)
		return
	}
	bt := map[string]interface{}{
		"code":        0,
		"key":         key,
		"image":       masterImageBase64,
		"tile":        tileImageBase64,
		"tile_width":  blockData.Width,
		"tile_height": blockData.Height,
		"tile_x":      blockData.TileX,
		"tile_y":      blockData.TileY,
	}
	result.OkWithData(bt, c)
}

// CheckSlideData 验证滑动验证码
// @Summary 验证滑动验证码
// @Description 验证滑动验证码
// @Tags 验证滑动验证码
// @Param point query string true "点击坐标"
// @Param key query string true "验证码key"
// @Success 200 {string} json
// @Router /api/captcha/shape/slide/check [get]
func (CaptchaAPI) CheckSlideData(c *gin.Context) {
	point := c.Query("point")
	key := c.Query("key")
	if point == "" || key == "" {
		result.FailWithNull(c)
		return
	}

	cacheDataByte, err := redis.Get(key).Bytes()
	if len(cacheDataByte) == 0 || err != nil {
		result.FailWithNull(c)
		return
	}
	src := strings.Split(point, ",")

	var dct *slide.Block
	if err := json.Unmarshal(cacheDataByte, &dct); err != nil {
		result.FailWithNull(c)
		return
	}

	chkRet := false
	if 2 == len(src) {
		sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[0]), 64)
		sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[1]), 64)
		chkRet = slide.CheckPoint(int64(sx), int64(sy), int64(dct.X), int64(dct.Y), 4)
	}

	if chkRet {
		result.OkWithMessage("success", c)
		return
	}
	result.FailWithMessage("fail", c)
}

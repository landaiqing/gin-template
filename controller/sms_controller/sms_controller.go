package sms_controller

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	gosms "github.com/pkg6/go-sms"
	"github.com/pkg6/go-sms/gateways"
	"github.com/pkg6/go-sms/gateways/aliyun"
	"github.com/pkg6/go-sms/gateways/smsbao"
	"schisandra-cloud-album/common/constant"
	"schisandra-cloud-album/common/redis"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
	"time"
)

type SmsController struct{}

// SendMessageByAli 发送短信验证码
// @Summary 发送短信验证码
// @Description 发送短信验证码
// @Tags 短信验证码
// @Produce json
// @Param phone query string true "手机号"
// @Router /controller/sms/ali/send [get]
func (SmsController) SendMessageByAli(c *gin.Context) {
	smsRequest := SmsRequest{}
	err := c.ShouldBindJSON(&smsRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	checkRotateData := utils.CheckRotateData(smsRequest.Angle, smsRequest.Key)
	if !checkRotateData {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
		return
	}
	isPhone := utils.IsPhone(smsRequest.Phone)
	if !isPhone {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}
	val := redis.Get(constant.UserLoginSmsRedisKey + smsRequest.Phone).Val()
	if val != "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaTooOften"), c)
		return
	}
	sms := gosms.NewParser(gateways.Gateways{
		ALiYun: aliyun.ALiYun{
			Host:            global.CONFIG.SMS.Ali.Host,
			AccessKeyId:     global.CONFIG.SMS.Ali.AccessKeyID,
			AccessKeySecret: global.CONFIG.SMS.Ali.AccessKeySecret,
		},
	})
	code := utils.GenValidateCode(6)
	wrong := redis.Set(constant.UserLoginSmsRedisKey+smsRequest.Phone, code, time.Minute).Err()
	if wrong != nil {
		global.LOG.Error(wrong)
		return
	}
	_, err = sms.Send(smsRequest.Phone, gosms.MapStringAny{
		"content":  "您的验证码是：****。请不要把验证码泄露给其他人。",
		"template": global.CONFIG.SMS.Ali.TemplateID,
		//"signName": global.CONFIG.SMS.Ali.Signature,
		"data": gosms.MapStrings{
			"code": code,
		},
	}, nil)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendSuccess"), c)

}

// SendMessageBySmsBao 短信宝发送短信验证码
// @Summary 短信宝发送短信验证码
// @Description 短信宝发送短信验证码
// @Tags 短信验证码
// @Produce json
// @Param phone query string true "手机号"
// @Router /controller/sms/smsbao/send [post]
func (SmsController) SendMessageBySmsBao(c *gin.Context) {
	smsRequest := SmsRequest{}
	err := c.ShouldBindJSON(&smsRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	checkRotateData := utils.CheckRotateData(smsRequest.Angle, smsRequest.Key)
	if !checkRotateData {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
		return
	}
	isPhone := utils.IsPhone(smsRequest.Phone)
	if !isPhone {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneErrorFormat"), c)
		return
	}
	val := redis.Get(constant.UserLoginSmsRedisKey + smsRequest.Phone).Val()
	if val != "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaTooOften"), c)
		return
	}
	sms := gosms.NewParser(gateways.Gateways{
		SmsBao: smsbao.SmsBao{
			User:     global.CONFIG.SMS.SmsBao.User,
			Password: global.CONFIG.SMS.SmsBao.Password,
		},
	})
	code := utils.GenValidateCode(6)
	wrong := redis.Set(constant.UserLoginSmsRedisKey+smsRequest.Phone, code, time.Minute).Err()
	if wrong != nil {
		global.LOG.Error(wrong)
		return
	}
	_, err = sms.Send(smsRequest.Phone, gosms.MapStringAny{
		"content": "您的验证码是：" + code + "。请不要把验证码泄露给其他人。",
	}, nil)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendSuccess"), c)
}

// SendMessageTest 发送测试短信验证码
// @Summary 发送测试短信验证码
// @Description 发送测试短信验证码
// @Tags 短信验证码
// @Produce json
// @Param phone query string true "手机号"
// @Router /controller/sms/test/send [post]
func (SmsController) SendMessageTest(c *gin.Context) {
	smsRequest := SmsRequest{}
	err := c.ShouldBindJSON(&smsRequest)
	if err != nil {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	checkRotateData := utils.CheckRotateData(smsRequest.Angle, smsRequest.Key)
	if !checkRotateData {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaVerifyError"), c)
		return
	}
	isPhone := utils.IsPhone(smsRequest.Phone)
	if !isPhone {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneError"), c)
		return
	}
	code := utils.GenValidateCode(6)
	err = redis.Set(constant.UserLoginSmsRedisKey+smsRequest.Phone, code, time.Minute).Err()
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendSuccess"), c)

}

package sms_api

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	gosms "github.com/pkg6/go-sms"
	"github.com/pkg6/go-sms/gateways"
	"github.com/pkg6/go-sms/gateways/aliyun"
	"github.com/pkg6/go-sms/gateways/smsbao"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"schisandra-cloud-album/utils"
)

// SendMessageByAli 发送短信验证码
// @Summary 发送短信验证码
// @Description 发送短信验证码
// @Tags 短信验证码
// @Produce json
// @Param phone query string true "手机号"
// @Router /api/sms/ali/send [get]
func (SmsAPI) SendMessageByAli(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotEmpty"), c)
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
	_, err := sms.Send(phone, gosms.MapStringAny{
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
}

// SendMessageBySmsBao 短信宝发送短信验证码
// @Summary 短信宝发送短信验证码
// @Description 发送短信验证码
// @Tags 短信验证码
// @Produce json
// @Param phone query string true "手机号"
// @Router /api/sms/smsbao/send [get]
func (SmsAPI) SendMessageBySmsBao(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		result.FailWithMessage(ginI18n.MustGetMessage(c, "PhoneNotEmpty"), c)
		return
	}
	sms := gosms.NewParser(gateways.Gateways{
		SmsBao: smsbao.SmsBao{
			User:     global.CONFIG.SMS.SmsBao.User,
			Password: global.CONFIG.SMS.SmsBao.Password,
		},
	})
	code := utils.GenValidateCode(6)
	_, err := sms.Send(phone, gosms.MapStringAny{
		"content": "您的验证码是：" + code + "。请不要把验证码泄露给其他人。",
	}, nil)
	if err != nil {
		global.LOG.Error(err)
		result.FailWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendFailed"), c)
		return
	}
	result.OkWithMessage(ginI18n.MustGetMessage(c, "CaptchaSendSuccess"), c)
}

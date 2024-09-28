package sms_api

type SmsRequest struct {
	Phone string `json:"phone" binding:"required"`
	Angle int64  `json:"angle" binding:"required"`
	Key   string `json:"key" binding:"required"`
}

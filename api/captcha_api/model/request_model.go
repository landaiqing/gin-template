package model

type RotateCaptchaRequest struct {
	Angle int    `json:"angle"`
	Key   string `json:"key"`
}

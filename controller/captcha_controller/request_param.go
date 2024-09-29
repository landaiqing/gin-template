package captcha_controller

type RotateCaptchaRequest struct {
	Angle int    `json:"angle" binding:"required"`
	Key   string `json:"key" binding:"required"`
}

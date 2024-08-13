package dto

import "encoding/json"

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone   string `json:"phone"`
	Captcha string `json:"captcha"`
}

// ResponseData 返回数据
type ResponseData struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresAt    int64   `json:"expires_at"`
	UID          *string `json:"uid"`
}

func (res ResponseData) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res ResponseData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &res)
}

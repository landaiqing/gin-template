package dto

import "encoding/json"

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone     string `json:"phone" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
	AutoLogin bool   `json:"auto_login" binding:"required"`
}

// AccountLoginRequest 账号登录请求
type AccountLoginRequest struct {
	Account   string `json:"account" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AutoLogin bool   `json:"auto_login" binding:"required"`
}

// AddUserRequest 新增用户请求
type AddUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Repassword string `json:"repassword" binding:"required"`
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

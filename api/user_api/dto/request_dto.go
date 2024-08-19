package dto

import "encoding/json"

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PhoneLoginRequest 手机号登录请求
type PhoneLoginRequest struct {
	Phone     string `json:"phone"`
	Captcha   string `json:"captcha"`
	AutoLogin bool   `json:"auto_login"`
}

// AccountLoginRequest 账号登录请求
type AccountLoginRequest struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	AutoLogin bool   `json:"auto_login"`
}

// AddUserRequest 新增用户请求
type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Phone      string `json:"phone"`
	Captcha    string `json:"captcha"`
	Password   string `json:"password"`
	Repassword string `json:"repassword"`
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

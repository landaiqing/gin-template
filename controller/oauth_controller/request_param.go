package oauth_controller

import (
	"encoding/json"
	"time"
)

// ResponseData 返回数据
type ResponseData struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	UID          *string  `json:"uid"`
	UserInfo     UserInfo `json:"user_info"`
}
type UserInfo struct {
	Username string    `json:"username,omitempty"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	Phone    string    `json:"phone,omitempty"`
	Email    string    `json:"email,omitempty"`
	Gender   string    `json:"gender"`
	Status   int64     `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

func (res ResponseData) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res ResponseData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &res)
}

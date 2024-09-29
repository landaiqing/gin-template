package oauth_controller

import "encoding/json"

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

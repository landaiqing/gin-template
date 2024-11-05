package types

import (
	"encoding/json"
)

// ResponseData 返回数据
type ResponseData struct {
	AccessToken string  `json:"access_token"`
	UID         *string `json:"uid"`
	Username    string  `json:"username,omitempty"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	Status      int64   `json:"status"`
}

func (res ResponseData) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res ResponseData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &res)
}

type RedisToken struct {
	AccessToken string `json:"access_token"`
	UID         string `json:"uid"`
}

func (res RedisToken) MarshalBinary() ([]byte, error) {
	return json.Marshal(res)
}

func (res RedisToken) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &res)
}

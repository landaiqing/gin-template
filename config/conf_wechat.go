package config

type Wechat struct {
	AppID     string `json:"app-id"`
	AppSecret string `json:"app-secret"`
	Token     string `json:"token"`
	AESKey    string `json:"aes-key"`
}

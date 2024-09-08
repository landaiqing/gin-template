package config

type Wechat struct {
	AppID     string `yaml:"app-id"`
	AppSecret string `yaml:"app-secret"`
	Token     string `yaml:"token"`
	AESKey    string `yaml:"aes-key"`
}

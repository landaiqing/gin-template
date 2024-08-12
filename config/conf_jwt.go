package config

type JWT struct {
	Secret            string `yaml:"secret"`
	Expiration        string `yaml:"expiration"`
	RefreshExpiration string `yaml:"refresh-expiration"`
	RefreshTokenKey   string `yaml:"refresh-token-key"`
	HeaderKey         string `yaml:"header-key"`
	HeaderPrefix      string `yaml:"header-prefix"`
}

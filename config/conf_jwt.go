package config

type JWT struct {
	Secret       string `yaml:"secret"`
	HeaderKey    string `yaml:"header-key"`
	HeaderPrefix string `yaml:"header-prefix"`
	Issuer       string `yaml:"issuer"`
}

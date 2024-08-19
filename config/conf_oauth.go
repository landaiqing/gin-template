package config

// OAuth is the configuration of OAuth.
type OAuth struct {
	Github Github `yaml:"github"`
	Gitee  Gitee  `yaml:"gitee"`
	QQ     QQ     `yaml:"qq"`
}

// Github and GiteeConfig are the configuration of Github and Gitee OAuth.
type Github struct {
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	RedirectURI  string `yaml:"redirect-uri"`
}

// Gitee is the configuration of Gitee OAuth.
type Gitee struct {
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	RedirectURI  string `yaml:"redirect-uri"`
}

// QQ is the configuration of QQ OAuth.
type QQ struct {
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	RedirectURI  string `yaml:"redirect-uri"`
}

package config

// OAuth is the configuration of OAuth.
type OAuth struct {
	Github Github `yaml:"github"`
	Gitee  Gitee  `yaml:"gitee"`
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

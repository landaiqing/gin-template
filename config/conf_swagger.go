package config

type Swagger struct {
	Enabled     bool   `yaml:"enable"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
}

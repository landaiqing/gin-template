package config

type Encrypt struct {
	Key        string `yaml:"key"`
	IV         string `yaml:"iv"`
	PublicKey  string `yaml:"public-key"`
	PrivateKey string `yaml:"private-key"`
}

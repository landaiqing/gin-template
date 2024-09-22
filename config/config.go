package config

type Config struct {
	MySQL   MySQL   `yaml:"mysql"`
	Logger  Logger  `yaml:"logger"`
	System  System  `yaml:"system"`
	Redis   Redis   `yaml:"redis"`
	SMS     SMS     `yaml:"sms"`
	JWT     JWT     `yaml:"jwt"`
	Encrypt Encrypt `yaml:"encrypt"`
	Wechat  Wechat  `yaml:"wechat"`
	OAuth   OAuth   `yaml:"oauth"`
	Swagger Swagger `yaml:"swagger"`
	Casbin  Casbin  `yaml:"casbin"`
	MongoDB MongoDB `yaml:"mongodb"`
}

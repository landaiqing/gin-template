package config

type Config struct {
	MySQL  MySQL  `yaml:"mysql"`
	Logger Logger `yaml:"logger"`
	System System `yaml:"system"`
	Redis  Redis  `yaml:"redis"`
	SMS    SMS    `yaml:"sms"`
	JWT    JWT    `yaml:"jwt"`
}

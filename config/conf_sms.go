package config

type SMS struct {
	Ali    Ali    `yaml:"ali"`     //阿里云短信配置
	SmsBao SmsBao `yaml:"sms-bao"` //短信宝配置
}

type Ali struct {
	Host            string `yaml:"host"`              //主机地址
	AccessKeyID     string `yaml:"access-key-id"`     //阿里云AccessKeyId
	AccessKeySecret string `yaml:"access-key-secret"` //阿里云AccessKeySecret
	TemplateID      string `yaml:"template-id"`       //短信模板ID
	Signature       string `yaml:"signature"`         //短信签名
}
type SmsBao struct {
	User     string `yaml:"user"`     //短信宝用户名
	Password string `yaml:"password"` //短信宝密码
}

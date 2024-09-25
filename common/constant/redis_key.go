package constant

const (
	// 登录相关的redis key
	UserLoginSmsRedisKey     = "user:sms:"
	UserLoginTokenRedisKey   = "user:token:"
	UserLoginCaptchaRedisKey = "user:captcha:"
	UserLoginClientRedisKey  = "user:client:"
	UserLoginQrcodeRedisKey  = "user:qrcode:"
	UserSessionRedisKey      = "user:session:"
)

// 登录之后
const (
	CommentSubmitCaptchaRedisKey = "comment:submit:captcha:"
)

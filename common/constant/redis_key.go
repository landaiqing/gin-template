package constant

// 登录相关的redis key
const (
	UserLoginSmsRedisKey     = "user:sms:"
	UserLoginTokenRedisKey   = "user:token:"
	UserLoginCaptchaRedisKey = "user:captcha:"
	UserLoginClientRedisKey  = "user:client:"
	UserLoginQrcodeRedisKey  = "user:qrcode:"
	UserSessionRedisKey      = "user:session:"
)

// 评论相关的redis key
const (
	CommentSubmitCaptchaRedisKey  = "comment:submit:captcha:"
	CommentOfflineMessageRedisKey = "comment:offline:message:"
	CommentLikeLockRedisKey       = "comment:like:lock:"
	CommentDislikeLockRedisKey    = "comment:dislike:lock:"
	CommentLikeListRedisKey       = "comment:like:list:"
	CommentUserListRedisKey       = "comment:user:list:"
)

// 系统相关的redis key

const (
	SystemApiNonceRedisKey = "system:controller:nonce:"
)

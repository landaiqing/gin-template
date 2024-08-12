package result

type ErrorCode int

const (
	SystemError          ErrorCode = 1001
	ParamsError          ErrorCode = 1002
	ParamsNullError      ErrorCode = 1003
	ParamsFormatError    ErrorCode = 1004
	ParamsValueError     ErrorCode = 1005
	ParamsRangeError     ErrorCode = 1006
	ParamsRepeatError    ErrorCode = 1007
	ParamsMatchError     ErrorCode = 1008
	ParamsNotUniqueError ErrorCode = 1009
	FileSizeExceeded     ErrorCode = 1010
	CaptchaExpireError   ErrorCode = 1011
)

type ErrorMap map[ErrorCode]string

var ErrMap = ErrorMap{
	1001: "系统错误",
	1002: "参数错误",
	1003: "参数为空",
	1004: "参数格式错误",
	1005: "参数值错误",
	1006: "参数值范围错误",
	1007: "参数值重复",
	1008: "参数值不匹配",
	1009: "参数值不唯一",
	1010: "超出文件上传大小限制",
	1011: "验证码已过期",
}

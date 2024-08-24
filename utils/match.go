package utils

import "regexp"

// IsEmail 判断是否为邮箱
func IsEmail(email string) bool {
	// 邮箱的正则表达式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

// IsPhone 判断是否为手机号
func IsPhone(phone string) bool {
	// 手机号的正则表达式，这里以中国大陆的手机号为例
	phoneRegex := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(phoneRegex, phone)
	return match
}

// IsUsername 用户名的正则表达式
func IsUsername(username string) bool {
	/**
	1.用户名仅能使用数字，大小写字母和下划线。
	2.用户名中的数字必须在最后。 数字可以有零个或多个。
	3.用户名不能以数字开头。 用户名字母可以是小写字母和大写字母。
	4.用户名长度必须至少为3个字符。 两位用户名只能使用字母，最多20个字符
	*/
	phoneRegex := `^[a-zA-Z_]{2,18}[0-9]*$`
	match, _ := regexp.MatchString(phoneRegex, username)
	return match
}

// IsPassword 密码的正则表达式
func IsPassword(password string) bool {
	phoneRegex := `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{6,18}$`
	match, _ := regexp.MatchString(phoneRegex, password)
	return match
}

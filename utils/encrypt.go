package utils

import "golang.org/x/crypto/bcrypt"

// Encrypt 加密
func Encrypt(val string) (string, error) {
	// 使用bcrypt库的GenerateFromPassword函数进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

// Verify 验证
func Verify(hashedVal string, val string) bool {
	// 使用bcrypt库的CompareHashAndPassword函数比较密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedVal), []byte(val))
	return err == nil
}

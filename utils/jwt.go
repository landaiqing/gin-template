package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wumansgy/goEncrypt/aes"
	"schisandra-cloud-album/global"
	"time"
)

type JWTPayload struct {
	UserID *string  `json:"user_id"`
	RoleID []*int64 `json:"role_id"`
}

type JWTClaims struct {
	JWTPayload
	jwt.RegisteredClaims
}

var MySecret []byte

// GenerateAccessToken generates a JWT token with the given payload
func GenerateAccessToken(payload JWTPayload) (string, error) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	claims := JWTClaims{
		JWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(MySecret)
	if err != nil {
		return "", err
	}
	accessToken, err := aes.AesCtrEncryptHex([]byte(signedString), []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return accessToken, nil
}

// GenerateAccessTokenAndRefreshToken generates a JWT token with the given payload, and returns the accessToken and refreshToken
func GenerateAccessTokenAndRefreshToken(payload JWTPayload) (string, string, int64) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	// accessToken 的数据
	accessClaims := JWTClaims{
		JWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.CONFIG.JWT.Issuer,
		},
	}
	refreshClaims := JWTClaims{
		JWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.CONFIG.JWT.Issuer,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	accessTokenString, err := accessToken.SignedString(MySecret)
	if err != nil {
		global.LOG.Error(err)
		return "", "", 0
	}
	refreshTokenString, err := refreshToken.SignedString(MySecret)
	if err != nil {
		global.LOG.Error(err)
		return "", "", 0
	}
	accessTokenEncrypted, err := aes.AesCtrEncryptHex([]byte(accessTokenString), []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	if err != nil {
		fmt.Println(err)
		return "", "", 0
	}
	refreshTokenEncrypted, err := aes.AesCtrEncryptHex([]byte(refreshTokenString), []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	if err != nil {
		fmt.Println(err)
		return "", "", 0
	}
	return accessTokenEncrypted, refreshTokenEncrypted, refreshClaims.ExpiresAt.Time.Unix()
}

// ParseToken parses a JWT token and returns the payload
func ParseToken(tokenString string) (*JWTPayload, bool, error) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	plaintext, err := aes.AesCtrDecryptByHex(tokenString, []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	if err != nil {
		global.LOG.Error(err)
		return nil, false, err
	}
	token, err := jwt.ParseWithClaims(string(plaintext), &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		global.LOG.Error(err)
		return nil, false, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &claims.JWTPayload, true, nil
	}
	return nil, false, err
}

package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"schisandra-cloud-album/global"
)

type RefreshJWTPayload struct {
	UserID *string `json:"user_id"`
	Type   *string `json:"type" default:"refresh"`
}
type AccessJWTPayload struct {
	UserID *string `json:"user_id"`
	Type   *string `json:"type" default:"access"`
}
type AccessJWTClaims struct {
	AccessJWTPayload
	jwt.RegisteredClaims
}
type RefreshJWTClaims struct {
	RefreshJWTPayload
	jwt.RegisteredClaims
}

var MySecret []byte

// GenerateAccessToken generates a JWT token with the given payload
func GenerateAccessToken(payload AccessJWTPayload) (string, error) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	claims := AccessJWTClaims{
		AccessJWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(MySecret)
	if err != nil {
		return "", err
	}
	// accessToken, err := aes.AesCtrEncryptHex([]byte(signedString), []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }
	return accessToken, nil
}

// GenerateRefreshToken generates a JWT token with the given payload, and returns the accessToken and refreshToken
func GenerateRefreshToken(payload RefreshJWTPayload, days time.Duration) string {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	refreshClaims := RefreshJWTClaims{
		RefreshJWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(days)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.CONFIG.JWT.Issuer,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(MySecret)
	if err != nil {
		global.LOG.Error(err)
		return ""
	}
	// refreshTokenEncrypted, err := aes.AesCtrEncryptHex([]byte(refreshTokenString), []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", 0
	// }
	return refreshTokenString
}

// ParseAccessToken parses a JWT token and returns the payload
func ParseAccessToken(tokenString string) (*AccessJWTPayload, bool, error) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	// plaintext, err := aes.AesCtrDecryptByHex(tokenString, []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	// if err != nil {
	// 	global.LOG.Error(err)
	// 	return nil, false, err
	// }
	token, err := jwt.ParseWithClaims(tokenString, &AccessJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, false, err
	}
	if claims, ok := token.Claims.(*AccessJWTClaims); ok && token.Valid {
		return &claims.AccessJWTPayload, true, nil
	}
	return nil, false, err
}

// ParseRefreshToken parses a JWT token and returns the payload
func ParseRefreshToken(tokenString string) (*RefreshJWTPayload, bool, error) {
	MySecret = []byte(global.CONFIG.JWT.Secret)
	// plaintext, err := aes.AesCtrDecryptByHex(tokenString, []byte(global.CONFIG.Encrypt.Key), []byte(global.CONFIG.Encrypt.IV))
	// if err != nil {
	// 	global.LOG.Error(err)
	// 	return nil, false, err
	// }
	token, err := jwt.ParseWithClaims(tokenString, &RefreshJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		global.LOG.Error(err)
		return nil, false, err
	}
	if claims, ok := token.Claims.(*RefreshJWTClaims); ok && token.Valid {
		return &claims.RefreshJWTPayload, true, nil
	}
	return nil, false, err
}

package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"schisandra-cloud-album/global"
	"time"
)

type JWTPayload struct {
	UserID   int    `json:"user_id"`
	Role     string `json:"role"`
	Username string `json:"username"`
}

type JWTClaims struct {
	JWTPayload
	jwt.RegisteredClaims
}

var MySecret = []byte(global.CONFIG.JWT.Secret)

// GenerateToken generates a JWT token with the given payload
func GenerateToken(payload JWTPayload) (string, error) {
	claims := JWTClaims{
		JWTPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(MySecret)
}

// ParseToken parses a JWT token and returns the payload
func ParseToken(tokenString string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		global.LOG.Error(err)
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &claims.JWTPayload, nil
	}
	return nil, err
}

package logic

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/huizm/go-chatroom/model"
	"time"
)

const (
	jwtSecret = "cohadsega"
	jwtExpire = 86400 * time.Second
)

type Claims struct {
	User *model.User `json:"user"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	return []byte(jwtSecret)
}

func GenerateToken(u *model.User) (string, error) {
	now := time.Now()
	expire := now.Add(jwtExpire)

	claims := &Claims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
			Issuer:    "go-chatroom",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if token != nil {
		claims, ok := token.Claims.(*Claims)
		if ok && token.Valid {
			return claims, nil
		}
	}

	return nil, err
}

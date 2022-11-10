package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shjting0510/sa_user/option"
	"time"
)

// TokenExpireDuration token有效期
const TokenExpireDuration = time.Hour * 24

var salt = []byte("4G%Bd7kaw3")

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	// token有效期
	tokenExpireDuration time.Duration
	jwt.RegisteredClaims
}

func NewCustomClaims(userID int64, opts ...option.Option[CustomClaims]) *CustomClaims {
	c := &CustomClaims{
		UserID:              userID,
		tokenExpireDuration: TokenExpireDuration,
	}

	option.Options[CustomClaims](opts).Apply(c)

	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.tokenExpireDuration)),
		Issuer:    "test",
	}
	return c
}

func WithUserID(userID int64) option.Option[CustomClaims] {
	return func(c *CustomClaims) {
		c.UserID = userID
	}
}

func WithExpire(expire time.Duration) option.Option[CustomClaims] {
	return func(c *CustomClaims) {
		c.tokenExpireDuration = expire
	}
}

// GenToken 生成JWT
func (claims *CustomClaims) GenToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(salt)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return salt, nil
		})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserToken struct {
	*jwt.StandardClaims
	Uid      uint
	Nickname string
}

// 密钥
const secret = "jehR3MRvMqPoQAA3oWdPbKKNuw0kmjtIC1aFhVzX6esfzIXCmoPQMzgCmF7KvcCPx9KWQH3khAjNXwJGex2x4FYMw17f2ZWXZEIq3tXtpNJHORnXUQkDjhAemOPhG5uH"

// 过期时间
const expires = time.Hour * 24 * 3650

// CreateJWT 创建JWT
func CreateJWT(uid uint, nickname string) string {
	nowTime := time.Now()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserToken{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: nowTime.Add(expires).Unix(), // 设置有效期
			Issuer:    "Owner",                     // 签发人
			IssuedAt:  nowTime.Unix(),              //签发时间
		},
		Uid:      uid,
		Nickname: nickname,
	})
	// 使用指定的secret签名,获取字符串token
	signedString, _ := claims.SignedString([]byte(secret))

	return signedString
}

// ValidJWT 解析JWT
func ValidJWT(jwtStr string) (UserToken, error) {
	token, err := jwt.ParseWithClaims(jwtStr, &UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return UserToken{}, err
	}

	// 断言类型
	claim, ok := token.Claims.(*UserToken)

	// 验证
	if !ok || !token.Valid {
		return UserToken{}, errors.New("token不合法")
	}

	return *claim, nil
}

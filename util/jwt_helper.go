// 用于生成和解析Token
package util

import (
	"time"

	"github.com/MrGaoRock666/Mairuida/user_service/config"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义Claims结构体，包含用户ID
type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成JWT Token
func GenerateJWT(userID uint64) (string, error) {
	// 构造声明
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpire)), // 过期时间
			Issuer:    config.Issuer,                                          // 签发者
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
		},
	}

	// 使用HS256算法生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名
	return token.SignedString(config.JwtSecret)
}

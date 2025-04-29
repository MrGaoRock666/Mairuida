// 用于配置JWT密钥
package config

import "time"

// JWT相关配置
var (
	JwtSecret   = []byte("MairuidaSuperSecretKey") // 密钥，TODO:真实项目中请放到环境变量或配置中心
	TokenExpire = time.Hour * 24 * 7               // Token 过期时间：7天
	Issuer      = "mairuida_user_service"          // 签发者
)

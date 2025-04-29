// 缓存连接
package config

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ct = context.Background()

// 初始化 Redis 连接
func InitRedis() *redis.Client {
	// 创建 Redis 客户端
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 没有密码
		DB:       0,                // 默认 DB
	})

	// 测试 Redis 连接
	_, err := RedisClient.Ping(ct).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	fmt.Println("Redis connected successfully!")
	return RedisClient
}

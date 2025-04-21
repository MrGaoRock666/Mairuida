// 数据库连接
package config

import (
	"Mairuida/user_service/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库连接
func InitDB() *gorm.DB {
	// MySQL 连接配置，替换为你实际的用户名、密码和数据库名
	dsn := "root:@Gaobaolin200510@tcp(127.0.0.1:3306)/mairuida_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移建表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return db
}

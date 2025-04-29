package config

import (
	"log"

	orderModel "github.com/MrGaoRock666/Mairuida/order_service/model"
	userModel "github.com/MrGaoRock666/Mairuida/user_service/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化 order_service 数据库连接
func InitOrderDB() *gorm.DB {
	// MySQL 连接配置
	// mysql.Open(dsn) 用于创建一个 MySQL 驱动实例
	dsn := "root:@Gaobaolin200510@tcp(127.0.0.1:3306)/mairuida_order?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect order database: %v", err)
	}

	// 自动迁移建表,根据 orderModel.Order 结构体的定义创建或更新数据库表
	// 若表已存在，会根据结构体定义更新表结构
	err = db.AutoMigrate(&orderModel.Order{})
	if err != nil {
		// 若自动迁移失败，记录错误日志并终止程序
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// 返回初始化的 gorm.DB 实例
	return db
}

// 初始化 user_service 数据库连接
func InitUserDB() *gorm.DB {
	dsn := "root:@Gaobaolin200510@tcp(127.0.0.1:3306)/mairuida_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect user database: %v", err)
	}

	// 自动迁移建表，根据 userModel.Address 结构体的定义创建或更新数据库表
	// 若表已存在，会根据结构体定义更新表结构
	err = db.AutoMigrate(&userModel.Address{})
	if err != nil {
		log.Fatalf("AutoMigrate for user database failed: %v", err)
	}

	return db
}

// 根据地址 ID 查询实际地址信息
func GetAddressByID(userDB *gorm.DB, addressID string) (userModel.Address, error) {
	var address userModel.Address
	result := userDB.Where("id = ?", addressID).First(&address)
	if result.Error != nil {
		return address, result.Error
	}
	return address, nil
}

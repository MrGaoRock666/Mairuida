package model

import (
	"time"

	"gorm.io/gorm"
)

// 订单状态的枚举定义，应与 proto 中保持一致
type OrderStatus int

const (
	StatusUnknown     OrderStatus = iota //未知状态
	StatusCreated                        //订单已创建
	StatusDispatching                    //订单待配送
	StatusInTransit                      //订单配送中
	StatusDelivered                      //订单已送达
	StatusCancelled                      //订单已取消
)

// 订单模型结构体
type Order struct {
	ID                    uint           `gorm:"primaryKey;autoIncrement" json:"id"`            // 数据库主键，自增ID
	OrderID               string         `gorm:"type:varchar(64);uniqueIndex" json:"order_id"`  // 订单ID，系统生成
	UserID                uint           `gorm:"not null" json:"user_id"`                       // 用户ID，外键字段
	SenderAddress         string         `gorm:"type:varchar(255)" json:"sender_address"`       // 发件地址
	ReceiverAddress       string         `gorm:"type:varchar(255)" json:"receiver_address"`     // 收件地址
	ItemName              string         `gorm:"type:varchar(64)" json:"item_name"`             // 物品名称
	Weight                float64        `gorm:"not null" json:"weight"`                        // 重量（kg）
	Volume                float64        `gorm:"not null" json:"volume"`                        // 体积（立方米）
	LogisticsCompany      string         `gorm:"type:varchar(128)" json:"logistics_company"`    // 指定物流公司名称
	PreferredDeliveryTime time.Time      `gorm:"type:timestamp" json:"preferred_delivery_time"` // 用户期望送达时间
	IsUrgent              bool           `gorm:"default:false" json:"is_urgent"`                // 是否加急
	IsDelayed             bool           `gorm:"default:false" json:"is_delayed"`               // 是否延迟
	Status                OrderStatus    `gorm:"default:1" json:"status"`                       // 状态（初始为 Created）
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`              // 创建时间
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" json:"updated_at"`              // 更新时间
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`                                // 软删除字段
}

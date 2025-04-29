package model

import (
	"time"

	"gorm.io/gorm"
)

// 根据protobuf接口定义出的待存储的字段
// 这里面有些字段是无需存储入数据库的
// 该结构体映射数据库的用户表
type User struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement"`         //数据库主键，自增ID，和用户ID其实可以是一个东西
	Username  string    `gorm:"type:varchar(64);unique;not null"` //用户名，不能为空，而且必须唯一
	Password  string    `gorm:"type:varchar(128);not null"`       //密码
	Email     string    `gorm:"type:varchar(128);unique"`         //邮箱
	Phone     string    `gorm:"type:varchar(20)"`                 //电话号码
	IsVIP     bool      //是否是VIP用户
	IsDeleted bool      `gorm:"default:false" json:"-"` // 是否注销：默认为 false，未注销
	VIPLevel  string    `gorm:"type:varchar(32)"`       //VIP等级
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	Addresses []Address `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // 一对多关系，一个用户可以对应多个地址
}

// Model文件夹用于存储Gorm模型结构体
package model

import (
	"gorm.io/gorm"
)

// 根据protobuf接口定义出待存储的字段
// 该结构体映射数据库的地址表
type Address struct {
	gorm.Model
	ID       string `gorm:"primaryKey;type:varchar(64)"`         // 数据库主键
	UserID   uint   `gorm:"not null;index;type:bigint unsigned"` // 修改为 unsigned bigint
	User     User   `gorm:"constraint:OnDelete:CASCADE;"`        // 外键约束
	Label    string `gorm:"type:varchar(64)"`                    //标签：如家，学校，公司等
	Province string `gorm:"type:varchar(64)"`                    //省
	City     string `gorm:"type:varchar(64)"`                    //市
	District string `gorm:"type:varchar(64)"`                    //区
	Detail   string `gorm:"type:text"`                           //细节用文本格式存储就行
}

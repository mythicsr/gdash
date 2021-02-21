package mysql_aaa

import (
	"time"
)

const (
	OriginTableName     = "aaa"
	TableNameDivideMode = "" //分表类型 day=按天 month=按月 默认不分表
)

//https://sql2gorm.mccode.info/
type AAA struct {
	Id         int64  `gorm:"column:id;primary_key"` // 主键
	Uid        int64  `gorm:"column:uid"`            // 用户id
	Info       string `gorm:"column:info;NOT NULL"`  // 用户信息
	InfoStruct string `gorm:"-"`                     // 用户信息Struct
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func GetCurrentTableName() string {
	cur := OriginTableName
	return cur
}

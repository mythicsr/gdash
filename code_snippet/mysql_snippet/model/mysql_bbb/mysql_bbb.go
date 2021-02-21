package mysql_bbb

import (
	"fmt"
	"github.com/jinzhu/now"
	"time"
)

const (
	OriginTableName     = "bbb"
	TableNameDivideMode = "day" //分表类型 day=按天 month=按月 默认不分表
)

//https://sql2gorm.mccode.info/
type BBB struct {
	Id        int64  `gorm:"column:id;primary_key"` // 主键
	Uid       int64  `gorm:"column:uid"`            // 用户id
	Info      string `gorm:"column:info;NOT NULL"`  // 用户信息
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetCurrentTableName() string {
	cur := OriginTableName

	switch TableNameDivideMode {
	case "day":
		{
			cur = fmt.Sprintf("%s%s", OriginTableName, time.Now().Format("20060102"))
		}
	case "month":
		{
			cur = fmt.Sprintf("%s%s", OriginTableName, time.Now().Format("200601"))
		}
	default:
		{

		}
	}
	return cur
}

func GetRecentTableNames(originTableName string, divideMode string, sTime, eTime time.Time) []string {
	var tables []string

	if divideMode == "day" {
		end := now.New(eTime).EndOfDay()
		for t := sTime; t.Before(end); {
			tb := fmt.Sprintf("%s%s", originTableName, t.Format("20060102"))
			tables = append(tables, tb)

			t = t.Add(24 * time.Hour)
		}
		return tables
	}

	if divideMode == "month" {
		end := now.New(eTime).EndOfMonth()
		for t := now.New(sTime); t.Before(end); {
			tb := fmt.Sprintf("%s%s", originTableName, t.Format("200601"))
			tables = append(tables, tb)

			t = now.New(t.EndOfMonth().Add(time.Second))
		}
		return tables
	}

	return []string{originTableName}
}

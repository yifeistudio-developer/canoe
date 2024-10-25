package db

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Name  string `gorm:"type:varchar(64);not null;comment:'姓名';default:''"`
	Type  int    `gorm:"type:integer;not null;comment:'类型';default:0"`
	RelId int64  `gorm:"type:bigint;not null;comment:'关联id';default:0"`
}

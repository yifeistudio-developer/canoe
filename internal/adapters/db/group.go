package db

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Attr         int    `gorm:"type:integer;not null;comment:'属性';default:0"`
	Name         string `gorm:"type:varchar(64);not null;comment:'名称';default:''"`
	Status       int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	Avatar       string `gorm:"type:varchar(256);not null;comment:'头像';default:''"`
	MemberCounts int    `gorm:"type:integer;not null;comment:'群成员数';default:0"`
}

package db

import "gorm.io/gorm"

type GroupMember struct {
	gorm.Model
	Attr   int    `gorm:"type:integer;not null;comment:'成员属性';default:0"`
	Status int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	Name   string `gorm:"type:varchar(64);not null;comment:'群昵称';default:''"`
	UserId int64  `gorm:"type:bigint;not null;comment:'用户id';default:0"`
}

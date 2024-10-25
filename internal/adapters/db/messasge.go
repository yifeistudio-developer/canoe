package db

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Type      int         `gorm:"type:integer;not null;comment:'类型';default:0"`
	Attr      int         `gorm:"type:integer;not null;comment:'属性';default:0"`
	SessionId int64       `gorm:"type:bigint;not null;comment:'会话id';default:0"`
	UserId    int64       `gorm:"type:bigint;not null;comment:'发送人id';default:0"`
	Payload   interface{} `gorm:"type:json;null;comment:'负载信息'"`
}

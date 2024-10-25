package db

import "gorm.io/gorm"

type UserSession struct {
	gorm.Model
	UserId    int64 `gorm:"type:bigint;not null;comment:'用户id';default:0"`
	SessionId int64 `gorm:"type:bigint;not null;comment:'会话id';default:0"`
	MsgCur    int64 `gorm:"type:bigint;not null;comment:'消息消费指针';default:0"`
}

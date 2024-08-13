package data

import "gorm.io/gorm"

// 用户信息

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(64);not null;comment:'姓名';default:''"`
	Avatar    string `gorm:"type:varchar(256);not null;comment:'头像';default:''"`
	Status    int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	AccountId int64  `gorm:"type:bigint;not null;comment:'关联账号id';default:0"`
}

func (User) TableName() string {
	return "t_user"
}

// 群组

type Group struct {
	gorm.Model
	Attr         int    `gorm:"type:integer;not null;comment:'属性';default:0"`
	Name         string `gorm:"type:varchar(64);not null;comment:'名称';default:''"`
	Status       int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	Avatar       string `gorm:"type:varchar(256);not null;comment:'头像';default:''"`
	MemberCounts int    `gorm:"type:integer;not null;comment:'群成员数';default:0"`
}

func (Group) TableName() string {
	return "t_group"
}

// 群成员

type GroupMember struct {
	gorm.Model
	Attr   int    `gorm:"type:integer;not null;comment:'成员属性';default:0"`
	Status int    `gorm:"type:integer;not null;comment:'状态';default:0"`
	Name   string `gorm:"type:varchar(64);not null;comment:'群昵称';default:''"`
	UserId int64  `gorm:"type:bigint;not null;comment:'用户id';default:0"`
}

func (GroupMember) TableName() string {
	return "t_group_member"
}

// 会话

type Session struct {
	gorm.Model
	Name  string `gorm:"type:varchar(64);not null;comment:'姓名';default:''"`
	Type  int    `gorm:"type:integer;not null;comment:'类型';default:0"`
	RelId int64  `gorm:"type:bigint;not null;comment:'关联id';default:0"`
}

func (Session) TableName() string {
	return "t_session"
}

// 用户会话

type UserSession struct {
	gorm.Model
	UserId    int64 `gorm:"type:bigint;not null;comment:'用户id';default:0"`
	SessionId int64 `gorm:"type:bigint;not null;comment:'会话id';default:0"`
	MsgCur    int64 `gorm:"type:bigint;not null;comment:'消息消费指针';default:0"`
}

func (UserSession) TableName() string {
	return "t_user_session"
}

// 消息

type Message struct {
	gorm.Model
	Type      int         `gorm:"type:integer;not null;comment:'类型';default:0"`
	Attr      int         `gorm:"type:integer;not null;comment:'属性';default:0"`
	SessionId int64       `gorm:"type:bigint;not null;comment:'会话id';default:0"`
	UserId    int64       `gorm:"type:bigint;not null;comment:'发送人id';default:0"`
	Payload   interface{} `gorm:"type:json;null;comment:'负载信息'"`
}

func (Message) TableName() string {
	return "t_message"
}

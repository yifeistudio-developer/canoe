package model

// 用户信息

type User struct {
	Name      string
	Avatar    string
	Status    int
	AccountId int64
}

// 群组

type Group struct {
	Attr         int
	Name         string
	Status       int
	Avatar       string
	MemberCounts int
}

// 群成员

type GroupMember struct {
	Attr   int
	Status int
	Name   string
	UserId int64
}

// 会话

type Session struct {
	Name  string
	Type  int
	RelId int64
}

// 用户会话

type UserSession struct {
	UserId    int64
	SessionId int64
	MsgCur    int64
}

// 消息

type Message struct {
	Type      int
	Attr      int
	SessionId int64
	UserId    int64
	Payload   interface{}
}

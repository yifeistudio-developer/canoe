package domain

type User struct {
	Name      string
	Avatar    string
	Status    int
	AccountId int64
}

type Group struct {
	Attr         int
	Name         string
	Status       int
	Avatar       string
	MemberCounts int
}

type GroupMember struct {
	Attr   int
	Status int
	Name   string
	UserId int64
}

type Session struct {
	Name  string
	Type  int
	RelId int64
}

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

type Envelope struct {
	Payload interface{}
}

type AlpsUserProfile struct {
	Username string
	Avatar   string
	Nickname string
}

type Result struct {
	Code      int         `json:"code"`
	IsSuccess bool        `json:"isSuccess"`
	Msg       string      `json:"msg,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Event struct {
}

func Success(data interface{}) *Result {
	return &Result{
		Code:      200,
		Data:      data,
		IsSuccess: true,
	}
}

func Fail(code int, msg string) Result {
	return Result{
		Code:      code,
		Msg:       msg,
		IsSuccess: false,
	}
}

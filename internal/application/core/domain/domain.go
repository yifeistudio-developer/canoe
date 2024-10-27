package domain

import (
	"encoding/json"
	"fmt"
)

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

func (r Result) Error() string {
	marshal, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("{\"code\":%d,\"is_success\":%v,\"msg\":\"%s\"}", r.Code, r.IsSuccess, r.Msg)
	}
	return string(marshal)
}

type Event struct {
}

func SuccessNil() Result {
	return Success(nil)
}

func Success(data interface{}) Result {
	return Result{
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

package model

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Result {
	return Result{
		Code: 0,
		Data: data,
	}
}

func Fail(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
	}
}

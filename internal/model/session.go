package model

type Session struct {

	// 会话ID
	Id int64 `json:"id"`

	// 单聊：0 群聊：1
	Type int `json:"type"`

	// 单聊：用户ID，群聊：群ID
	RelId int64 `json:"relId"`
}

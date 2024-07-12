package model

type Message struct {
	From    int64       `json:"from"` // from session ID
	To      int64       `json:"to"`   // session ID
	Payload interface{} `json:"payload"`
}

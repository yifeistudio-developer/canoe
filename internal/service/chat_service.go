package service

import "github.com/kataras/golog"

var logger *golog.Logger

type ChatService struct {
}

func NewChatService() *ChatService {
	return &ChatService{}
}

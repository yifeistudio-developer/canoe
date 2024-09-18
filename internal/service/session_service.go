package service

import "canoe/internal/model"

type SessionService struct {
}

func NewSessionService() *SessionService {
	return &SessionService{}
}

func (*SessionService) NewSession(accessToken string) *model.Session {
	// 根据accessToken 获取用户基本信息
	return &model.Session{}
}

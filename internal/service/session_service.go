package service

import (
	"canoe/internal/model"
	"canoe/internal/model/data"
)

type SessionService struct {
}

func NewSessionService() *SessionService {
	return &SessionService{}
}

func (*SessionService) NewSession(accessToken string) *model.Session {
	// 根据accessToken 获取用户基本信息
	return &model.Session{}
}

func (*SessionService) ListSession(username string, cur int) []model.Session {
	if cur <= 0 {
		cur = 1
	}
	size := 10
	var sessions []data.Session
	db.Where("username = ?", username).
		Offset((cur - 1) * size).
		Limit(size).Find(&sessions)
	result := make([]model.Session, len(sessions))
	for i, s := range sessions {
		result[i] = model.Session{
			Name: s.Name,
		}
	}
	return result
}

package api

import (
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

type UserService struct {
	db ports.UserDbPort
}

func (u UserService) Register(user domain.User) error {
	return nil
}

package api

import (
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

type UserApi struct {
	db ports.UserDbPort
}

func (u UserApi) Register(user domain.User) error {
	return nil
}

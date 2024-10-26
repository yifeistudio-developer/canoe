package ports

import "github.com/yifeistudio-developer/canoe/internal/application/core/domain"

type ApiPort interface {
	GetUserApiPort() UserApiPort
}

type UserApiPort interface {
	Register(user domain.User) error
}

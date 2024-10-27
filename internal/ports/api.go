package ports

import "github.com/yifeistudio-developer/canoe/internal/application/core/domain"

type ApiPort interface {
	UserApiPort() UserApiPort
}

type UserApiPort interface {
	Register(user domain.User) error
}

package ports

import "github.com/yifeistudio-developer/canoe/internal/application/core/domain"

type UserApiPort interface {
	Register(user domain.User) error
}

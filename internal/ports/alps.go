package ports

import "github.com/yifeistudio-developer/canoe/internal/application/core/domain"

type AlpsAdapter interface {
	GetAccountPrincipals() (domain.AlpsUserProfile, error)
}

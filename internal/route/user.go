package route

import (
	. "canoe/internal/model"
	"canoe/internal/service"
)

type userController struct {
	Us *service.UserService
}

func (c *userController) Get() *Result {
	return Success(nil)
}

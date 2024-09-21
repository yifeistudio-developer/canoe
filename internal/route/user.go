package route

import (
	. "canoe/internal/model"
	"canoe/internal/service"
)

type userController struct {
	userService *service.UserService
}

func (c *userController) Get() *Result {
	return Success(nil)
}

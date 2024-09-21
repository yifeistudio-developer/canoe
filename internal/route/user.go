package route

import (
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

var (
	userService = service.NewUserService()
)

type userController struct {
}

func UserRoutes(route router.Party) {

	route.Get("/{id:int64}", func(ctx iris.Context) {

	})

	route.Post("/{id:int64}", func(ctx iris.Context) {

	})

	route.Put("/{id:int64}", func(ctx iris.Context) {

	})

	route.Delete("/{id:int64}", func(ctx iris.Context) {

	})
}

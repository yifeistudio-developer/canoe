package route

import (
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func UserRoutes(route router.Party, userService *service.UserService) {

	route.Get("/{id:int64}", func(ctx iris.Context) {

	})

	route.Post("/{id:int64}", func(ctx iris.Context) {

	})

	route.Put("/{id:int64}", func(ctx iris.Context) {

	})

	route.Delete("/{id:int64}", func(ctx iris.Context) {

	})
}

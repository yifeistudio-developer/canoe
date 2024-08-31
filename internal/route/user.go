package route

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

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

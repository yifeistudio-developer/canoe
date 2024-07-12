package main

import (
	"canoe/internal/config"
	"canoe/internal/route"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	cfg := config.LoadConfig()
	app.Logger().Install(config.GetLogger(cfg.Logger))
	app.UseError(config.GlobalErrorHandler)
	app.UseGlobal(config.AccessLogger)
	route.SetupRoutes(app)
	err := app.Listen(cfg.Server.Address)
	if err != nil {
		app.Logger().Error("failed to start server: ", err.Error())
		panic("failed to start server: " + err.Error())
	}
}

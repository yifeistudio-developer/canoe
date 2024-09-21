package main

import (
	"canoe/internal/config"
	"canoe/internal/route"
	"github.com/kataras/iris/v12"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"strconv"
)

func startIris(cfg *config.Config) (*iris.Application, bool) {
	app := iris.Default()
	app.Logger().Install(config.GetLogger(cfg.Logger))
	app.UseGlobal(config.AccessLogger)
	app.UseError(config.GlobalErrorHandler)
	route.SetupRoutes(app)
	signal := make(chan bool)
	defer close(signal)
	go func() {
		err := app.Listen(":"+strconv.Itoa(int(cfg.Server.Port)), func(application *iris.Application) {
			//
			signal <- true
		})
		if err != nil {
			logger.Error("failed to start server: ", err.Error())
			signal <- false
		}
	}()
	return app, <-signal
}

package web

import (
	"github.com/kataras/iris/v12"
	"github.com/yifeistudio-developer/canoe/config"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/middleware"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/route"
	"github.com/yifeistudio-developer/canoe/internal/ports"
)

func NewAdapter(app ports.ApiPort) *iris.Application {
	server := iris.Default()
	accessLog := middleware.NewAccessLog(config.GetLogPath())
	server.UseGlobal(middleware.ContextErrorHandler)
	server.UseRouter(accessLog.Handler)
	server.UseError(middleware.ErrorHandler)
	party := server.Party("/canoe/api")
	route.Register(party, app)
	return server
}

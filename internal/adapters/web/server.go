package web

import (
	"github.com/kataras/iris/v12"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/middleware"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web/route"
	"github.com/yifeistudio-developer/canoe/internal/application/core/api"
	"strconv"
)

type Adapter struct {
	port   int
	app    *api.Application
	server *iris.Application
}

func NewAdapter(port int, app *api.Application) *Adapter {
	return &Adapter{port: port, app: app}
}

func (a Adapter) Startup(logPath string) {
	server := iris.Default()
	accessLog := middleware.NewAccessLog(logPath)
	server.UseRouter(accessLog.Handler)
	party := server.Party("/canoe/api")
	route.Register(party, a.app)
	s := make(chan bool)
	defer close(s)
	go func() {
		err := server.Listen(":"+strconv.Itoa(a.port), func(application *iris.Application) {
			s <- true
		})
		if err != nil {
			server.Logger().Error("failed to start server: ", err.Error())
			s <- false
		}
	}()
	<-s
	a.server = server
}

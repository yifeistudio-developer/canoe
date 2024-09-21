package route

import (
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
)

type sessionController struct {
	sessionService *service.SessionService
}

func (s *sessionController) listSession(ctx iris.Context) {

}

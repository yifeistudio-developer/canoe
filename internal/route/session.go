package route

import (
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
)

type sessionController struct {
	Ss *service.SessionService
}

func (s *sessionController) Get(ctx iris.Context) *Result {
	ss := s.Ss
	sessions := ss.ListSession("develop@yifeistudio.com")
	return Success(sessions)
}

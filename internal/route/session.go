package route

import (
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"strconv"
)

type sessionController struct {
	Ss *service.SessionService
}

func (s *sessionController) Get(ctx iris.Context) *Result {
	ss := s.Ss
	curStr := ctx.URLParam("cur")
	cur, err := strconv.Atoi(curStr)
	if err != nil {
		cur = 1
	}
	logger.Infof("cur: %v", cur)
	sessions := ss.ListSession("develop@yifeistudio.com", cur)
	return Success(sessions)
}

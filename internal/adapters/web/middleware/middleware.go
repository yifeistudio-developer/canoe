package middleware

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/yifeistudio-developer/canoe/internal/application/core/domain"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

func NewSocketServer() *WebsocketServer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &WebsocketServer{
		peers:  &sync.Map{},
		ctx:    ctx,
		cancel: cancelFunc,
	}
}

func NewAccessLog(logPath string) *accesslog.AccessLog {
	logFile := "access.log"
	fullPath := filepath.Join(logPath, logFile)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			log.Fatal("cannot create log directory")
		}
	}
	ac := accesslog.File(fullPath)
	ac.AddOutput(os.Stdout)
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = false
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler
	return ac
}

func ErrorHandler(ctx iris.Context) {
	code := ctx.GetStatusCode()
	err := ctx.StopWithJSON(code, domain.Fail(code, http.StatusText(code)))
	if err != nil {
	}
}

func ContextErrorHandler(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger := ctx.Application().Logger()
			logger.Error("handle error: path = ", ctx.Path(), " error = ", err)
			// unformed error forward to ErrorHandler
			if reflect.TypeOf(err) != reflect.TypeOf(domain.Result{}) {
				ctx.StatusCode(http.StatusInternalServerError)
				return
			}
			//goland:noinspection GoTypeAssertionOnErrors
			result := err.(domain.Result)
			// handle build-in error code
			if code := result.Code; http.StatusText(code) != "" {
				if result.Msg == "" {
					result.Msg = http.StatusText(code)
				}
				err = ctx.StopWithJSON(code, result)
				return
			}
			err = ctx.StopWithJSON(http.StatusOK, err.(domain.Result))
		}
	}()
	ctx.Next()
}

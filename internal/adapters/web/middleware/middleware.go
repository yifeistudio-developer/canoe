package middleware

import (
	"github.com/kataras/iris/v12/middleware/accesslog"
	"log"
	"os"
	"path/filepath"
)

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

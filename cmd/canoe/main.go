package main

import (
	"canoe/internal/config"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

var cfg *config.Config
var app *iris.Application
var rpc *grpc.Server

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	startup()
	<-quit
	shutdown()
}

// 启动
func startup() {
	cfg = config.LoadConfig()
	db := connectDB(cfg.Database)
	service.SetDB(db)
	isStarted := false
	app, isStarted = startIris(cfg)
	if !isStarted {
		panic("Canoe Web-Server is not started")
	}
	log := app.Logger()
	log.Info("Iris Server started")
	rpc = startGrpcServer(log)
	log.Info("Grpc Server started")
	// do register
	config.Register(cfg, log)
}

// 关闭
func shutdown() {
	log := app.Logger()
	// do de-register
	log.Info("Server is shutting down...")
	config.DeRegister(cfg, log)
	rpc.GracefulStop()
	log.Info("Grpc Server gracefully stopped")
	err := app.Shutdown(context.Background())
	if err == nil {
		log.Info("Iris Server gracefully stopped")
	}
	log.Info("Server exited")
	os.Exit(0)
}

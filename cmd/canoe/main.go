package main

import (
	"canoe/internal/config"
	"canoe/internal/database"
	"database/sql"
	"github.com/kataras/iris/v12"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	running := runningContext{}
	running.startup()
	<-quit
	running.shutdown()
}

type runningContext struct {
	cfg  *config.Config
	app  *iris.Application
	grpc *grpc.Server
	db   *sql.DB
}

// 启动
func (c *runningContext) startup() {
	cfg := config.LoadConfig()
	app, isStarted := startIris(cfg)
	if !isStarted {
		panic("Canoe Web-Server is not started")
	}
	log := app.Logger()
	log.Info("Iris Server started")
	gRpcServer := startGrpcServer(log)
	log.Info("Grpc Server started")
	// do register
	config.Register(cfg, log)
	c.cfg = cfg
	c.app = app
	c.grpc = gRpcServer
	c.db = database.Connect(cfg.Database)
	log.Info("Database connected")
}

// 关闭
func (c *runningContext) shutdown() {
	log := c.app.Logger()
	// do de-register
	log.Info("Server is shutting down...")
	config.DeRegister(c.cfg, log)
	c.grpc.GracefulStop()
	log.Info("Grpc Server gracefully stopped")
	err := c.app.Shutdown(context.Background())
	if err == nil {
		log.Info("Iris Server gracefully stopped")
	}
	database.Disconnect(c.db)
	log.Info("database disconnected")
	log.Info("Server exited")
	os.Exit(0)
}

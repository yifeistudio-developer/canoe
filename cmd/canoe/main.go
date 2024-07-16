package main

import (
	"canoe/internal/config"
	"canoe/internal/route"
	"github.com/kataras/iris/v12"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	WebServerStarted = iota
	GrpcServerStarted
)

func main() {
	app := iris.New()
	cfg := config.LoadConfig()
	app.Logger().Install(config.GetLogger(cfg.Logger))
	app.UseGlobal(config.AccessLogger)
	app.UseError(config.GlobalErrorHandler)
	route.SetupRoutes(app)
	var wg sync.WaitGroup
	logger := app.Logger()

	go func() {
		wg.Add(1)
		err := app.Listen(":"+strconv.Itoa(int(cfg.Server.Port)), func(application *iris.Application) {
			wg.Done()
		})
		if err != nil {
			logger.Error("failed to start server: ", err.Error())
			panic("failed to start server: " + err.Error())
		}
	}()

	go func() {
		wg.Add(1)
		time.Sleep(3000 * time.Millisecond)
		logger.Printf("starting grpc server at port %d", cfg.Server.Port)
		wg.Done()
	}()

	wg.Wait()
	// do register
	config.Register(cfg, logger)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// do de-register
	logger.Println("Server is shutting down...")
	config.DeRegister(cfg, logger)
	logger.Println("Server exited")
	os.Exit(0)
}

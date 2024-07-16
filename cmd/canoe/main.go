package main

import (
	"canoe/internal/config"
	remote "canoe/internal/grpc"
	"canoe/internal/grpc/generated"
	"canoe/internal/route"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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
	wg.Add(1)
	go func() {
		err := app.Listen(":"+strconv.Itoa(int(cfg.Server.Port)), func(application *iris.Application) {
			wg.Done()
		})
		if err != nil {
			logger.Error("failed to start server: ", err.Error())
			panic("failed to start server: " + err.Error())
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		generated.RegisterGreeterServer(s, &remote.Server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
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

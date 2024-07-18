package main

import (
	"canoe/internal/config"
	"canoe/internal/remote"
	"canoe/internal/remote/generated"
	"canoe/internal/route"
	"context"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
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
	interceptor := remote.ServerInterceptor{Logger: logger}
	gRpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor))
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		generated.RegisterAuthenticationServiceServer(gRpcServer, &remote.Server{})
		for k, v := range gRpcServer.GetServiceInfo() {
			logger.Info("service info: ", k, v)
		}

		if err := gRpcServer.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()
	wg.Wait()
	// do register
	config.Register(cfg, logger)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// do de-register
	logger.Info("Server is shutting down...")
	config.DeRegister(cfg, logger)
	gRpcServer.GracefulStop()
	logger.Info("Grpc Server gracefully stopped")
	err := app.Shutdown(context.Background())
	if err == nil {
		logger.Info("Iris Server gracefully stopped")
	}
	logger.Info("Server exited")
	os.Exit(0)
}

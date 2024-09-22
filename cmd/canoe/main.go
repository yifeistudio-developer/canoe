package main

import (
	"canoe/internal/config"
	"canoe/internal/model/data"
	"canoe/internal/remote"
	"canoe/internal/remote/generated"
	"canoe/internal/route"
	"canoe/internal/service"
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var cfg *config.Config
var app *iris.Application
var rpc *grpc.Server
var logger *golog.Logger

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	startup()
	<-quit
	shutdown()
}

// 启动
func startup() {
	// load configuration
	cfg = config.LoadConfig()

	// init application
	app = iris.Default()
	app.Logger().SetLevel("info")
	app.Logger().Install(config.GetLogger(cfg.Logger))
	app.UseGlobal(config.AccessLogger)
	app.UseError(config.GlobalErrorHandler)
	logger = app.Logger()

	// start grpc server
	rpc = startGrpcServer()
	logger.Info("Grpc Server started")
	// connect database
	db := connectDB(cfg.Database)
	service.SetupServices(db, logger)
	config.Register(cfg, logger, false)

	// start http server
	isStarted := false
	app, isStarted = startIris(cfg)
	if !isStarted {
		panic("Canoe Web-Server is not started")
	}
	logger.Info("Iris Server started")
	// already to serve then do register
	config.Register(cfg, logger, true)
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

func connectDB(cfg struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     uint64 `envconfig:"DB_PORT" default:"5432"`
	Username string `envconfig:"DB_USERNAME" default:"canoe"`
	Password string `envconfig:"DB_PASSWORD" default:"canoe110930008"`
	DbName   string `envconfig:"DB_NAME" default:"canoe"`
}) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbName, cfg.Username, cfg.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// migrate
	err = db.AutoMigrate(&data.User{})
	err = db.AutoMigrate(&data.Group{})
	err = db.AutoMigrate(&data.GroupMember{})
	err = db.AutoMigrate(&data.Session{})
	err = db.AutoMigrate(&data.UserSession{})
	err = db.AutoMigrate(&data.Message{})
	if err != nil {
		return nil
	}
	return db
}

func startGrpcServer() *grpc.Server {
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
	return gRpcServer
}

func startIris(cfg *config.Config) (*iris.Application, bool) {
	// config router
	route.SetupRoutes(app)
	s := make(chan bool)
	defer close(s)
	go func() {
		err := app.Listen(":"+strconv.Itoa(int(cfg.Server.Port)), func(application *iris.Application) {
			//
			s <- true
		})
		if err != nil {
			app.Logger().Error("failed to start server: ", err.Error())
			s <- false
		}
	}()
	return app, <-s
}

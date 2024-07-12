package config

import (
	"github.com/joho/godotenv"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Address string `envconfig:"SERVER_ADDRESS" default:":8080"`
	}
	Logger struct {
		Level string `envconfig:"LOG_LEVEL" default:"info"`
	}
}

func LoadConfig() *Config {
	var cfg Config
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	log.Printf("loading %s configuration", env)
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found", envFile)
	}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return &cfg
}

func GetLogger(cnf struct {
	Level string `envconfig:"LOG_LEVEL" default:"info"`
}) *golog.Logger {
	logger := golog.New()
	logger.SetLevel(cnf.Level)
	return logger
}

// AccessLogger access logger
func AccessLogger(ctx iris.Context) {
	start := time.Now()
	ctx.Next()
	duration := time.Since(start)
	ctx.Application().Logger().Printf("%s %s took %v", ctx.Method(), ctx.Path(), duration.Milliseconds())
}

func GlobalErrorHandler(ctx iris.Context) {
	println("hanlde error")
}

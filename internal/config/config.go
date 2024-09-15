package config

import (
	. "canoe/internal/model"
	"github.com/joho/godotenv"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kelseyhightower/envconfig"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"os"
	"time"
)

var sc []constant.ServerConfig
var cc constant.ClientConfig

type Config struct {
	ApplicationName string `envconfig:"APPLICATION_NAME" default:"canoe"`
	Server          struct {
		Port uint64 `envconfig:"SERVER_PORT" default:"8080"`
	}
	Nacos struct {
		Host     string `envconfig:"NACOS_HOST" default:"localhost"`
		Port     uint64 `envconfig:"NACOS_PORT" default:"8848"`
		Username string `envconfig:"NACOS_USERNAME" default:"nacos"`
		Password string `envconfig:"NACOS_PASSWORD" default:"nacos"`
	}
	Logger struct {
		Level string `envconfig:"LOG_LEVEL" default:"info"`
	}
	Database struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     uint64 `envconfig:"DB_PORT" default:"5432"`
		Username string `envconfig:"DB_USERNAME" default:"canoe"`
		Password string `envconfig:"DB_PASSWORD" default:"canoe110930008"`
		DbName   string `envconfig:"DB_NAME" default:"canoe"`
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
	nacos := cfg.Nacos
	sc = []constant.ServerConfig{
		*constant.NewServerConfig(nacos.Host, nacos.Port),
	}
	cc = *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
		constant.WithUsername(nacos.Username),
		constant.WithPassword(nacos.Password),
	)

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	_, err = client.GetConfig(vo.ConfigParam{
		Group:  "DEFAULT",
		DataId: "canoe",
		OnChange: func(namespace, group, dataId, data string) {

		},
	})
	if err != nil {
		//log.Fatalf("Failed to load configuration: %s", err)
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
	code := ctx.GetStatusCode()
	result := Fail(code, "")
	_ = ctx.JSON(result)
}

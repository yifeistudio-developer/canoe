package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV", "development")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL", "postgres://canoe:canoe110930008@localhost:5432/canoe?sslmode=disable&TimeZone=Asia/Shanghai")
}

func GetLogPath() string {
	return getEnvironmentValue("LOG_PATH", "./tmp")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT", "3000")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}
	return port
}

func GetAlpsServiceUrl() string {
	return getEnvironmentValue("ALPS_SERVICE_URL", "")
}

func getEnvironmentValue(key string, option string) string {
	if os.Getenv(key) == "" && option == "" {
		log.Fatalf("environment variable %s not set", key)
	}
	if os.Getenv(key) == "" {
		return option
	} else {
		return os.Getenv(key)
	}
}

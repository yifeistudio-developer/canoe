package main

import (
	"github.com/yifeistudio-developer/canoe/config"
	"github.com/yifeistudio-developer/canoe/internal/adapters/db"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web"
	"github.com/yifeistudio-developer/canoe/internal/application/core/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	webAdapter := web.NewAdapter(config.GetApplicationPort(), application)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	webAdapter.Startup(config.GetLogPath())
	<-quit
	// shutdown
	webAdapter.Shutdown()
	os.Exit(0)
}

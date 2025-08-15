package main

import (
	"fmt"
	"log"

	"github.com/yifeistudio-developer/canoe/config"
	"github.com/yifeistudio-developer/canoe/internal/adapters/db"
	"github.com/yifeistudio-developer/canoe/internal/adapters/web"
	"github.com/yifeistudio-developer/canoe/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	app := api.NewApplication(dbAdapter)
	webAdapter := web.NewAdapter(app)
	err = webAdapter.Listen(fmt.Sprintf(":%d", config.GetApplicationPort()))
	if err != nil {
		log.Fatalf("failed to start web server: %v", err)
	}
}

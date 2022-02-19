package main

import (
	"github.com/vvkosty/go_sprint_1/internal/app"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

func main() {
	var appConfig config.ServerConfig
	var appHandler handler.Handler

	application := app.App{
		Config:  &appConfig,
		Storage: storage.NewStorage(),
		Handler: &appHandler,
	}
	application.Init()
	application.Start()
}

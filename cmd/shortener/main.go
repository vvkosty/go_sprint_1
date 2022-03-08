package main

import (
	"github.com/vvkosty/go_sprint_1/internal/app"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	middleware "github.com/vvkosty/go_sprint_1/internal/app/middlewares"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

func main() {
	var appConfig config.ServerConfig
	var appHandler handler.Handler
	var appMiddleware middleware.Middleware

	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:     &appConfig,
		Handler:    &appHandler,
		Middleware: &appMiddleware,
	}

	if appConfig.FileStoragePath != "" {
		application.Storage = storage.NewFileStorage(appConfig.FileStoragePath)
	} else {
		application.Storage = storage.NewMapStorage()
	}
	defer application.Storage.Close()

	application.Init()
	application.Start()
}

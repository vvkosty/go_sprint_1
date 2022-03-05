package app

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

type App struct {
	Config  *config.ServerConfig
	Storage storage.Repository
	Handler *handler.Handler
}

func (app *App) Init() {
	app.Handler.Storage = app.Storage
	app.Handler.Config = app.Config
}

func (app *App) Start() {
	r := app.SetupRouter()

	err := r.Run(app.Config.Address)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *App) SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	r.GET("/:id", app.Handler.GetFullLink)
	r.POST("/", app.Handler.CreateShortLink)
	r.POST("/api/shorten", app.Handler.CreateJSONShortLink)

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	return r
}

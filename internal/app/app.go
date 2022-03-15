package app

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	middleware "github.com/vvkosty/go_sprint_1/internal/app/middlewares"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

type App struct {
	Config     *config.ServerConfig
	Storage    storage.Repository
	Handler    *handler.Handler
	Middleware *middleware.Middleware
}

func (app *App) Init() {
	app.Handler.Storage = app.Storage
	app.Handler.Config = app.Config
	app.Middleware.Config = app.Config
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
	r.Use(app.Middleware.SetCookie)

	r.GET("/:id", app.Handler.GetFullLink)
	r.POST("/", app.Handler.CreateShortLink)
	r.POST("/api/shorten", app.Handler.CreateJSONShortLink)
	r.GET("/api/user/urls", app.Handler.GetAllLinks)
	r.GET("/ping", app.Handler.Ping)
	r.POST("/api/shorten/batch", app.Handler.CreateBatchLinks)

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	return r
}

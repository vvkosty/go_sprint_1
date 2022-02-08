package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/config"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/handlers"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
)

func SetupRouter(urlStorage *handlers.Urls) *gin.Engine {
	r := gin.Default()
	r.GET("/:id", urlStorage.GetHandler)
	r.POST("/", urlStorage.PostHandler)

	r.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})

	return r
}

func Start() {
	urls := &handlers.Urls{DB: storage.NewMapDatabase()}
	r := SetupRouter(urls)

	err := r.Run(":" + config.ServerPort)
	if err != nil {
		fmt.Println(err)
	}
}

package server

import (
	"fmt"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/config"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/handlers"
	"net/http"
)

func Start() {
	http.HandleFunc("/", handlers.RootHandler)

	if err := http.ListenAndServe(":"+config.ServerPort, nil); err != nil {
		fmt.Println(err)
	}
}

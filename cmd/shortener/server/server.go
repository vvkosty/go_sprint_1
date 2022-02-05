package server

import (
	"fmt"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/config"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/handlers"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
	"net/http"
)

func Start() {
	urls := &handlers.Urls{DB: storage.NewMapDatabase()}

	http.HandleFunc("/", urls.RootHandler)

	if err := http.ListenAndServe(":"+config.ServerPort, nil); err != nil {
		fmt.Println(err)
	}
}

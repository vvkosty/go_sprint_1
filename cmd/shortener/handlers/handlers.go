package handlers

import (
	"fmt"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/config"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
	"io"
	"net/http"
	"net/url"
)

type Urls struct {
	DB storage.Repository
}

func (urls *Urls) RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		urls.GetHandler(w, r)
	case http.MethodPost:
		urls.PostHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (urls *Urls) GetHandler(w http.ResponseWriter, r *http.Request) {
	urlId := r.URL.Path[1:]
	originalUrl := urls.DB.Find(urlId)

	if len(originalUrl) > 0 {
		w.Header().Add("Location", originalUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (urls *Urls) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	urlToEncode, err := url.ParseRequestURI(string(body))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set(`Content-Type`, `plain/text`)
	checksum := urls.DB.Save(urlToEncode.String())

	fmt.Fprintf(w, "%s://%s:%s/%s", config.ServerScheme, config.ServerDomain, config.ServerPort, checksum)
}

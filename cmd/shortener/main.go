package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ServerScheme = "http"
	ServerDomain = "localhost"
	ServerPort   = "8080"
)

var DB = map[string]string{}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(w, r)
	case http.MethodPost:
		PostHandler(w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	urlId := r.URL.Path[1:]

	if originalUrl, ok := DB[urlId]; ok {
		w.Header().Add("Location", originalUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	urlToEncode, err := url.ParseRequestURI(string(body))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	checksum := strconv.Itoa(int(crc32.ChecksumIEEE([]byte(urlToEncode.String()))))
	DB[checksum] = urlToEncode.String()
	fmt.Fprintf(w, "%s://%s:%s/%s", ServerScheme, ServerDomain, ServerPort, checksum)
}

func main() {
	http.HandleFunc("/", RootHandler)

	if err := http.ListenAndServe(":"+ServerPort, nil); err != nil {
		fmt.Println(err)
	}
}

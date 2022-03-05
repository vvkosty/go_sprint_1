package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vvkosty/go_sprint_1/internal/app"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

func TestUrls_CreateShortLink(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "OK",
			request: `http://example.com/test-url/test1/test2/test.php`,
			want: want{
				code:     http.StatusCreated,
				response: `http://localhost:8080/3744865384`,
			},
		},
		{
			name:    "Empty url",
			request: ``,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
		{
			name:    "Incorrect url",
			request: `test/example.php`,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
	}

	application := createApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := application.SetupRouter()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func TestUrls_GetFullLink(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name  string
		urlID string
		want  want
	}{
		{
			name:  "OK",
			urlID: `3744865384`,
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: `http://example.com/test-url/test1/test2/test.php`,
			},
		},
		{
			name:  "Empty url",
			urlID: ``,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
		{
			name:  "Incorrect url",
			urlID: `test/example.php`,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
	}

	application := createApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application.Storage.Save("http://example.com/test-url/test1/test2/test.php")
			router := application.SetupRouter()
			request := httptest.NewRequest(http.MethodGet, "/"+tt.urlID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, res.Header.Get("Location"), tt.want.response)
		})
	}
}

func TestUrls_CreateJsonShortLink(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "OK",
			request: `{"url":"http://example.com/test-url/test1/test2/test.php"}`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/3744865384"}`,
			},
		},
		{
			name:    "Empty url",
			request: ``,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
		{
			name:    "Incorrect url",
			request: `test/example.php`,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
	}

	application := createApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := application.SetupRouter()
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.request))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func createApp() *app.App {
	var appConfig config.ServerConfig
	var appHandler handler.Handler

	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:  &appConfig,
		Storage: storage.NewMapStorage(),
		Handler: &appHandler,
	}
	application.Init()

	return &application
}

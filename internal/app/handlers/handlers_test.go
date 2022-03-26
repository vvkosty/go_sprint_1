package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vvkosty/go_sprint_1/internal/app"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	handler "github.com/vvkosty/go_sprint_1/internal/app/handlers"
	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
	middleware "github.com/vvkosty/go_sprint_1/internal/app/middlewares"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

var gzipHelper = helpers.GzipHelper{}

var appConfig config.ServerConfig
var appHandler handler.Handler
var appMiddleware middleware.Middleware

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
			application.Storage.Save(
				"http://example.com/test-url/test1/test2/test.php",
				"test123",
			)
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

func TestUrls_CreateBatchLinks(t *testing.T) {
	type (
		requestBatchURL struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		responseBatchURL struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}

		want struct {
			code     int
			response []responseBatchURL
		}
	)

	tests := []struct {
		name    string
		request []requestBatchURL
		want    want
	}{
		{
			name: "OK",
			request: []requestBatchURL{
				{
					CorrelationID: "123",
					OriginalURL:   "http://example.com/test-url/test1/test2/test.php",
				},
				{
					CorrelationID: "4576",
					OriginalURL:   "http://example.com/test-url/test1/test2/test1.php",
				},
			},
			want: want{
				code: http.StatusCreated,
				response: []responseBatchURL{
					{
						CorrelationID: "123",
						ShortURL:      "http://localhost:8080/3744865384",
					},
					{
						CorrelationID: "4576",
						ShortURL:      "http://localhost:8080/2566738836",
					},
				},
			},
		},
	}

	var responseBatchURLs []responseBatchURL
	application := createApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := application.SetupRouter()
			requestParams, _ := json.Marshal(tt.request)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(requestParams))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			json.Unmarshal(resBody, &responseBatchURLs)
			require.Equal(t, tt.want.response, responseBatchURLs)
			responseBatchURLs = responseBatchURLs[:0]
		})
	}
}

func TestUrls_CheckGZIPHeaders(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name                 string
		requestBody          string
		isCompressRequest    bool
		isDecompressResponse bool
		want                 want
	}{
		{
			name:                 "OK with full compress",
			requestBody:          `{"url":"http://example.com/test-url/test1/test2/test.php"}`,
			isCompressRequest:    true,
			isDecompressResponse: true,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"http://localhost:8080/3744865384"}`,
			},
		},
	}

	application := createApp()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := application.SetupRouter()

			requestBody := tt.requestBody
			if tt.isCompressRequest {
				compressedRequest, _ := gzipHelper.Compress([]byte(tt.requestBody))
				requestBody = string(compressedRequest)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(requestBody))

			if tt.isCompressRequest {
				request.Header.Set("Content-Encoding", "gzip")
			}
			if tt.isDecompressResponse {
				request.Header.Set("Accept-Encoding", "gzip")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			responseBody, _ := io.ReadAll(res.Body)
			if tt.isDecompressResponse {
				responseBody, _ = gzipHelper.Decompress(responseBody)
			}

			require.Equal(t, tt.want.response, string(responseBody))
		})
	}
}

func createApp() *app.App {
	appConfig.LoadEnvs()
	appConfig.ParseCommandLine()

	application := app.App{
		Config:     &appConfig,
		Storage:    storage.NewMapStorage(),
		Handler:    &appHandler,
		Middleware: &appMiddleware,
	}
	application.Init()

	return &application
}

func recreateTables(app *app.App) {
	app.Storage = storage.NewPostgresDatabase(appConfig.DatabaseDsn, true)
}

package handlers

import (
	"github.com/stretchr/testify/require"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUrls_PostHandler(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.request))

			urls := &Urls{DB: storage.NewMapDatabase()}

			w := httptest.NewRecorder()
			h := http.HandlerFunc(urls.PostHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func TestUrls_GetHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name  string
		urlId string
		want  want
	}{
		{
			name:  "OK",
			urlId: `3744865384`,
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: `http://example.com/test-url/test1/test2/test.php`,
			},
		},
		{
			name:  "Empty url",
			urlId: ``,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
		{
			name:  "Incorrect url",
			urlId: `test/example.php`,
			want: want{
				code:     http.StatusBadRequest,
				response: ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+tt.urlId, nil)

			db := storage.NewMapDatabase()
			db.Save("http://example.com/test-url/test1/test2/test.php")
			urls := &Urls{DB: db}

			w := httptest.NewRecorder()
			h := http.HandlerFunc(urls.GetHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.want.code, res.StatusCode)

			require.Equal(t, res.Header.Get("Location"), tt.want.response)
		})
	}
}

func TestUrls_RootHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		code   int
	}{
		{
			name:   "PUT",
			method: http.MethodPut,
			code:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", nil)

			urls := &Urls{DB: storage.NewMapDatabase()}

			w := httptest.NewRecorder()
			h := http.HandlerFunc(urls.RootHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			require.Equal(t, tt.code, res.StatusCode)
		})
	}
}

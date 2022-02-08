package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/handlers"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/server"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
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
			urls := &handlers.Urls{DB: storage.NewMapDatabase()}
			router := server.SetupRouter(urls)
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

func TestUrls_GetHandler(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := storage.NewMapDatabase()
			db.Save("http://example.com/test-url/test1/test2/test.php")
			urls := &handlers.Urls{DB: db}
			router := server.SetupRouter(urls)
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

package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
	"github.com/vvkosty/go_sprint_1/internal/app/helpers"
	storage "github.com/vvkosty/go_sprint_1/internal/app/storage"
)

type (
	Handler struct {
		Storage storage.Repository
		Config  *config.ServerConfig
	}

	requestURL struct {
		URL string `json:"url"`
	}

	responseURL struct {
		Result string `json:"result"`
	}

	listURL struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	requestBatchURL struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	responseBatchURL struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)

var err error

func (h *Handler) GetFullLink(c *gin.Context) {
	var entityDeletedError *storage.EntityDeletedError
	urlID := c.Param("id")
	originalURL, err := h.Storage.Find(urlID)
	if err != nil {
		if errors.As(err, &entityDeletedError) {
			c.Status(http.StatusGone)
			return
		}
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if len(originalURL) <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header(`Location`, originalURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (h *Handler) CreateShortLink(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	urlToEncode, err := url.ParseRequestURI(string(body))
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	checksum := helpers.GenerateChecksum(urlToEncode.String())
	entity, _ := h.Storage.Find(checksum)

	if entity != "" {
		c.Status(http.StatusConflict)
	} else {
		userID, _ := c.Get("userId")
		checksum, err = h.Storage.Save(urlToEncode.String(), userID.(string))
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusCreated)
	}

	c.Header(`Content-Type`, `plain/text`)
	responseBody := fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum)
	c.Writer.Write([]byte(responseBody))
}

func (h *Handler) CreateJSONShortLink(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	requestURL := requestURL{}
	if err := json.Unmarshal(body, &requestURL); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	checksum := helpers.GenerateChecksum(requestURL.URL)
	entity, _ := h.Storage.Find(checksum)

	c.Header(`Content-Type`, gin.MIMEJSON)
	if entity != "" {
		c.Status(http.StatusConflict)
	} else {
		userID, _ := c.Get("userId")
		checksum, err = h.Storage.Save(requestURL.URL, userID.(string))
		if err != nil {
			log.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusCreated)
	}

	response := responseURL{
		Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Writer.Write(encodedResponse)
}

func (h *Handler) GetAllLinks(c *gin.Context) {
	var response []listURL
	userID, _ := c.Get("userId")

	for checksum, originalURL := range h.Storage.List(userID.(string)) {
		response = append(response, listURL{
			ShortURL:    fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
			OriginalURL: originalURL,
		})
	}

	if len(response) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header(`Content-Type`, gin.MIMEJSON)
	c.Writer.Write(encodedResponse)
}

func (h *Handler) Ping(c *gin.Context) {
	db, err := sql.Open("pgx", h.Config.DatabaseDsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
}

func (h *Handler) CreateBatchLinks(c *gin.Context) {
	var requestBatchURLs []requestBatchURL
	var responseBatchURLs []responseBatchURL
	var uniqueViolatesError *storage.UniqueViolatesError

	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err := json.Unmarshal(body, &requestBatchURLs); err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header(`Content-Type`, gin.MIMEJSON)

	for _, requestURL := range requestBatchURLs {
		userID, _ := c.Get("userId")
		checksum, err := h.Storage.Save(requestURL.OriginalURL, userID.(string))
		if err != nil {
			log.Println(err)
			if errors.As(err, &uniqueViolatesError) {
				c.Status(http.StatusConflict)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}

		responseBatchURLs = append(responseBatchURLs, responseBatchURL{
			ShortURL:      fmt.Sprintf("%s/%s", h.Config.BaseURL, checksum),
			CorrelationID: requestURL.CorrelationID,
		})
	}

	encodedResponse, err := json.Marshal(&responseBatchURLs)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
	c.Writer.Write(encodedResponse)
}

func (h *Handler) DeleteBatchLinks(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	var checksums []string
	json.Unmarshal(body, &checksums)

	go h.Storage.DeleteBatchByChecksums(checksums)
	c.Status(http.StatusAccepted)
}

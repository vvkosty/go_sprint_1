package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	config "github.com/vvkosty/go_sprint_1/internal/app/config"
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
)

func (h *Handler) GetFullLink(c *gin.Context) {
	urlID := c.Param("id")
	originalURL := h.Storage.Find(urlID)

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

	c.Status(http.StatusCreated)
	c.Header(`Content-Type`, `plain/text`)
	checksum := h.Storage.Save(urlToEncode.String())
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

	c.Status(http.StatusCreated)
	c.Header(`Content-Type`, gin.MIMEJSON)
	checksum := h.Storage.Save(requestURL.URL)

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

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/config"
	"github.com/vvkosty/go_sprint_1/cmd/shortener/storage"
)

type (
	Urls struct {
		DB storage.Repository
	}

	requestURL struct {
		URL string `json:"url"`
	}

	responseURL struct {
		Result string `json:"result"`
	}
)

func (urls *Urls) GetFullLink(c *gin.Context) {
	urlID := c.Param("id")
	originalURL := urls.DB.Find(urlID)

	if len(originalURL) <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header(`Location`, originalURL)
	c.Status(http.StatusTemporaryRedirect)
}

func (urls *Urls) CreateShortLink(c *gin.Context) {
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
	checksum := urls.DB.Save(urlToEncode.String())
	responseBody := fmt.Sprintf("%s://%s:%s/%s", config.ServerScheme, config.ServerDomain, config.ServerPort, checksum)

	c.Writer.Write([]byte(responseBody))
}

func (urls *Urls) CreateJSONShortLink(c *gin.Context) {
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
	checksum := urls.DB.Save(requestURL.URL)

	response := responseURL{
		Result: fmt.Sprintf("%s://%s:%s/%s", config.ServerScheme, config.ServerDomain, config.ServerPort, checksum),
	}

	encodedResponse, err := json.Marshal(&response)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.Writer.Write(encodedResponse)
}

package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
)

type ShortURL struct {
	ShortURL string `json:"result"`
}

func Shorten(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать тело запроса"})
		return
	}

	URL := string(body)
	parsedURL, err := url.Parse(URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не спарсил URL"})
		return
	}

	shortURL := reduceURL()

	fmt.Printf("lin2klink2link---")
	fmt.Printf("Парсированный URL: %s\n", parsedURL.String())

	URLMap[shortURL] = URLPair{parsedURL, shortURL}

	result := ShortURL{ShortURL: Cfg.BaseURL + "/" + shortURL}
	fmt.Printf("Здесь json-----")
	c.JSON(http.StatusCreated, result)
}

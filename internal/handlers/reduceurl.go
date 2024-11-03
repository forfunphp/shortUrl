package handlers

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"shortUrl/config"
)

type URLPair struct {
	URL      *url.URL
	ShortURL string
}

type ShortURL struct {
	ShortURL string `json:"result"`
}

var URLMap = make(map[string]URLPair)
var Cfg = config.NewConfig()

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	err := Cfg.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}
}

func ReduceURL(c *gin.Context) {

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

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "text/html" {
		c.Data(http.StatusCreated, "text/html", []byte(Cfg.BaseURL+"/"+shortURL))
	} else if contentType == "text/plain" {
		c.Data(http.StatusCreated, "text/plain", []byte(Cfg.BaseURL+"/"+shortURL))
	} else if contentType == "application/json" {
		result := ShortURL{ShortURL: Cfg.BaseURL + "/" + shortURL}
		c.JSON(http.StatusCreated, result)
	}

}

func reduceURL() string {
	var shortURL string
	for i := 0; i < 8; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		shortURL += string(charset[randomIndex.Int64()])
	}
	return shortURL
}

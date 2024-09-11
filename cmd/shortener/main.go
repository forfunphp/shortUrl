package main

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
	"strings"
)

type URLPair struct {
	URL      *url.URL
	ShortURL string
}

var urlMap = make(map[string]URLPair)
var cfg = config.NewConfig()

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	err := cfg.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
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

func main() {
	router := gin.Default()

	router.POST("/", reduceURLHandler)
	router.GET("/:shortURL", redirectHandler)

	fmt.Printf("Сервер запущен на %s\n", cfg.HTTPAddr)
	colonIndex := strings.Index(cfg.HTTPAddr, ":")
	port := cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
}

func reduceURLHandler(c *gin.Context) {

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
	urlMap[shortURL] = URLPair{parsedURL, shortURL}

	c.Data(http.StatusCreated, "text/plain", []byte(cfg.BaseURL+shortURL))
}

func redirectHandler(c *gin.Context) {
	shortURL := c.Param("shortURL")
	urlPair, ok := urlMap[shortURL]

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет урла"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, urlPair.URL.String())
}

package handlers

import (
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		// Разархивируем данные
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось раз архивировать данные"})
			return
		}
		defer reader.Close()

		// Читаем тело запроса
		body, err = io.ReadAll(reader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать тело запроса"})
			return
		}
	}

	URL := string(body)
	parsedURL, err := url.Parse(URL)

	logger2, _ := zap.NewDevelopment()
	defer logger2.Sync()
	logger2.Info("Request xxxxx",
		zap.String("fullURL", URL),
		zap.String("fullURL", parsedURL.String()), // Добавляем полный URL
		zap.Int("Status", c.Writer.Status()),      // Добавляем parsedURL
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{URL + "error": "Не спарсил URL"})
		return
	}

	shortURL := reduceURL()

	fmt.Printf("lin2klink2link--->")
	fmt.Printf("Парсированный URL: %s\n", parsedURL.String())

	URLMap[shortURL] = URLPair{parsedURL, shortURL}

	contentType := c.Request.Header.Get("Content-Type")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info("Request processed1ss",
		zap.String("fullURL", c.Request.URL.String()), // Добавляем полный URL
		zap.String("parsedURL", parsedURL.String()),
		zap.String("contentType", contentType),
		zap.Int("Status", c.Writer.Status()), // Добавляем parsedURL
	)

	if contentType == "text/html" {
		c.Data(http.StatusCreated, "text/html", []byte(Cfg.BaseURL+"/"+shortURL))
	} else if contentType == "text/plain; charset=utf-8" {
		c.Data(http.StatusCreated, "text/plain", []byte(Cfg.BaseURL+"/"+shortURL))
	} else if contentType == "application/json" {

		result := ShortURL{ShortURL: Cfg.BaseURL + "/" + shortURL}
		c.JSON(http.StatusCreated, result)
	} else if contentType == "application/x-gzip" {

		// Отправляем сжатый ответ с Content-Type: application/x-gzip
		c.Data(http.StatusCreated, "application/x-gzip", []byte(Cfg.BaseURL+"/"+shortURL))
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

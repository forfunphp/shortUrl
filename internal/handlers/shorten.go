package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"shortUrl/config"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

var Cfg = config.NewConfig()

func Shorten(c *gin.Context) {

	log.Println("11111111122222223333333")

	if Cfg.Databes != "" {
		log.Println("9999444444444")
	}

	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}
	// Парсим URL из структуры запроса
	parsedURL, err := url.Parse(req.URL) // Используем req.Url

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось спарсить URL"})
		return
	}

	shortURL := reduceURL()

	URLMap[shortURL] = URLPair{parsedURL, shortURL}

	var resp ShortURL                            // Инициализируем структуру ShortURL
	resp.ShortURL = Cfg.BaseURL + "/" + shortURL // Заполняем поле ShortURL

	// Кодируем ответ в JSON

	jsonData, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "text/html" {
		c.Data(http.StatusCreated, "text/html; charset=utf-8", jsonData)
	} else if contentType == "text/plain; charset=utf-8" {
		c.Data(http.StatusCreated, "text/plain; charset=utf-8", jsonData)
	} else if contentType == "application/json" {
		//c.JSON(http.StatusCreated, result)
		c.Data(http.StatusCreated, "application/json", jsonData)
	}

	//c.Data(http.StatusCreated, "application/json", jsonData) // Удаляем string(jsonData)

}

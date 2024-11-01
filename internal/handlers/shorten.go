package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type ShortURL struct {
	ShortURL string `json:"result"`
}

type ShortenRequest struct {
	Url string `json:"url"`
}

func Shorten(c *gin.Context) {

	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	// Парсим URL из структуры запроса
	parsedURL, err := url.Parse(req.Url) // Используем req.Url
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось спарсить URL"})
		return
	}

	shortURL := reduceURL()

	fmt.Printf("Здесь json--->")
	fmt.Printf("Парсированный URL: %s\n", parsedURL.String())

	URLMap[shortURL] = URLPair{parsedURL, shortURL}

	var resp ShortURL                            // Инициализируем структуру ShortURL
	resp.ShortURL = Cfg.BaseURL + "/" + shortURL // Заполняем поле ShortURL

	// Кодируем ответ в JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем ответ
	c.JSON(http.StatusCreated, string(jsonData))

	fmt.Printf("<----Здесь json")

}

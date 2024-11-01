package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}
type ShortURL struct {
	ShortURL string `json:"result"` // Имя поля `result` соответствует JSON
}

func Shorten(c *gin.Context) {

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
	fmt.Println(string(jsonData))
	// Отправляем ответ
	c.Data(http.StatusCreated, "application/json", jsonData) // Удаляем string(jsonData)

}

package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
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

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info("Request processed1",
		zap.String("fullURL", c.Request.URL.String()), // Добавляем полный URL
		zap.String("parsedURL", string(jsonData)),     // Добавляем parsedURL
	)
	//fmt.Println(string(jsonData))

	c.Data(http.StatusInternalServerError, "application/x-gzip", jsonData)

}

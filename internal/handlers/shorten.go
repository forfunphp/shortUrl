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

	//fmt.Println(string(jsonData))
	contentType := c.Request.Header.Get("Content-Type")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	logger.Info("Request processed1",
		zap.String("fullURL", c.Request.URL.String()), // Добавляем полный URL
		zap.String("parsedURL", string(jsonData)),
		zap.String("contentType", contentType), // Добавляем parsedURL
		zap.Int("statusCode", c.Writer.Status()),
	)

	if contentType == "text/html" {
		c.Data(http.StatusCreated, "text/html; charset=utf-8", jsonData)
	} else if contentType == "text/plain; charset=utf-8" {
		c.Data(http.StatusCreated, "text/plain; charset=utf-8", jsonData)
	} else if contentType == "application/json" {
		//c.JSON(http.StatusCreated, string(jsonData))
		logger.Info("Request processed2",
			zap.String("fullURL2", c.Request.URL.String()), // Добавляем полный URL
			zap.String("parsedURL2", string(jsonData)),
			zap.String("contentType2", contentType), // Добавляем parsedURL
			zap.Int("statusCode2", c.Writer.Status()),
		)
		c.JSON(http.StatusCreated, jsonData)
	} else if contentType == "application/x-gzip" {
		c.Data(http.StatusCreated, "application/x-gzip", jsonData)
	}
	// Отправляем ответ
	//c.Data(http.StatusCreated, "application/json", jsonData) // Удаляем string(jsonData)

}

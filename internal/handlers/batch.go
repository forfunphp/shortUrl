package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
)

type BatchRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func Batch(c *gin.Context) {

	var batchRequest []BatchRequestEntry

	if err := c.BindJSON(&batchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	var batchResponse []BatchResponseEntry

	for _, reqEntry := range batchRequest {

		parsedURL, err := url.Parse(reqEntry.OriginalURL)
		if err != nil {

			batchResponse = append(batchResponse, BatchResponseEntry{
				CorrelationID: reqEntry.CorrelationID,
				ShortURL:      "ERROR: Не удалось спарсить URL",
			})
			log.Printf("Ошибка парсинга URL %s: %v", reqEntry.OriginalURL, err)
			continue
		}

		shortURL := reduceURL()

		URLMap[shortURL] = URLPair{parsedURL, shortURL}

		fullShortURL := Cfg.BaseURL + "/" + shortURL
		batchResponse = append(batchResponse, BatchResponseEntry{
			CorrelationID: reqEntry.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}

	c.JSON(http.StatusCreated, batchResponse)
}

package handlers

import (
	"encoding/json"
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

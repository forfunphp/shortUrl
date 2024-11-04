package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	if contentType == "text/html" {
		c.Data(http.StatusCreated, "text/html; charset=utf-8", jsonData)
	} else if contentType == "text/plain; charset=utf-8" {
		c.Data(http.StatusCreated, "text/plain; charset=utf-8", jsonData)
	} else if contentType == "application/json"
	{
		//c.JSON(http.StatusCreated, result)
		c.Data(http.StatusCreated, "application/json", jsonData)
	} else if contentType == "application/x-gzip"	{

		jsonBytes, err := json.Marshal(gin.H{"ShortURL": Cfg.BaseURL + "/" + shortURL})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сериализации JSON"})
			return
		}

		// Сжимаем JSON-ответ с помощью gzip
		var buffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&buffer)
		_, err = gzipWriter.Write(jsonBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сжатия gzip"})
			return
		}
		gzipWriter.Close()

		// Отправляем сжатый ответ с Content-Type: application/x-gzip
		c.Data(http.StatusCreated, "application/x-gzip", buffer.Bytes())
	}
	// Отправляем ответ
	//c.Data(http.StatusCreated, "application/json", jsonData) // Удаляем string(jsonData)

}

package handlers

import (
	"compress/gzip"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"shortUrl/config"
	"strings"
)

type URLPair struct {
	URL      *url.URL
	ShortURL string
}

type URLData struct {
	UUID        uuid.UUID `json:"uuid"` // Тип данных uuid.UUID
	ShortURL    string
	OriginalURL *url.URL
}

type ShortURL struct {
	ShortURL string `json:"result"`
}

type OriginalURL struct {
	ShortURL string `json:"original_url"`
}

var URLMap = make(map[string]URLPair)
var URLDat = make(map[string]URLData)
var Cfg = config.NewConfig()

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	err := Cfg.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}
}

func ReduceURL(c *gin.Context) {

	if Cfg.Databes != "" {
		log.Println("99999999999999999999999999999999")
	}

	log.Println("5666666666666")

	body, err := readRequestBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать тело запроса"})
		return
	}

	// Разбор URL
	parsedURL, err := url.Parse(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный URL"})
		return
	}

	shortURL := reduceURL()

	fmt.Printf("lin2klink2link--->")
	fmt.Printf("Парсированный URL: %s\n", parsedURL.String())

	URLMap[shortURL] = URLPair{parsedURL, shortURL}

	var urls []URLData
	urls = append(urls, URLData{
		UUID:        uuid.New(),
		ShortURL:    shortURL,
		OriginalURL: parsedURL,
	})

	if Cfg.Databes != "" {

		db, err := sql.Open("postgres", Cfg.Databes) // Замените "postgres" именем вашего драйвера
		if err != nil {
			log.Printf("не удалось открыть базу данных: %v", err)
		}

		_, err = db.Exec("INSERT INTO short_urls (id, shortURL, parsedURL) VALUES ($1, $2, $3)", uuid.New(), shortURL, parsedURL)
		if err != nil {
			log.Println("111111111111111111")
			log.Printf("Error saving to database: %v", err)
		}
		defer db.Close()
	}

	filePath := Cfg.EnvFilePath
	saveURLsToFile(urls, filePath)
	contentType := c.Request.Header.Get("Content-Type")

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

func saveURLsToFile(urls []URLData, fname string) error {
	data, err := json.MarshalIndent(urls, "", " ")
	if err != nil {
		return err
	}

	// Проверка существования файла
	_, err = os.Stat(fname)
	if os.IsNotExist(err) {
		// Создание файла, если он не существует
		file, err := os.Create(fname)
		if err != nil {
			return fmt.Errorf("не удалось создать файл %s: %w", fname, err)
		}
		defer file.Close()

		logger5, _ := zap.NewDevelopment()
		defer logger5.Sync()

		logger5.Info("Request processed33ffddd33",
			zap.String("fname", fname),
		)

		return os.WriteFile(fname, data, 0666)
	} else if err != nil {
		return fmt.Errorf("не удалось получить доступ к файлу %s: %w", fname, err)
	}

	// Файл существует, записываем данные в него
	return os.WriteFile(fname, data, 0666)
}

func readRequestBody(c *gin.Context) ([]byte, error) {
	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)
	} else {
		return io.ReadAll(c.Request.Body)
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

package main

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"net/url"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"os"
	"shortUrl/config"
	"shortUrl/internal/handlers"
	"strings"
	"time"
)

var logger *zap.Logger
var db *sql.DB
var sugar zap.SugaredLogger
var Cfg = config.NewConfig()

type URLData struct {
	UUID        uuid.UUID `json:"uuid"` // Тип данных uuid.UUID
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

func main() {

	dsn := os.Getenv("DATABASE_DSN")

	params, err := parsePostgresDSN(dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Parsed parameters:", params)
		log.Fatal(err)
	}

	if dsn != "" {
		handlers.NewPostgresStore(dsn)
	}

	filePath := Cfg.EnvFilePath
	loadURLsFromFile(filePath)

	router := gin.Default()
	router.Use(gzipMiddleware())

	router.POST("/", WithLogging(handlers.ReduceURL))
	router.GET("/:shortURL", WithLogging(handlers.Redirect))
	router.POST("/api/shorten", WithLogging(handlers.Shorten))
	router.GET("/ping", WithLogging(handlers.Ping))

	//_, err := NewPostgresStore(dsn)

	fmt.Printf("Сервер запущен на %s\n", handlers.Cfg.HTTPAddr)

	colonIndex := strings.Index(handlers.Cfg.HTTPAddr, ":")
	if colonIndex == -1 {
		log.Fatalf("Неверный формат адреса: %s", handlers.Cfg.HTTPAddr)
	}

	port := handlers.Cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
}

func parsePostgresDSN(dsn string) (map[string]string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %w", err)
	}

	params := make(map[string]string)
	params["host"] = u.Hostname()
	params["port"] = u.Port()
	params["user"] = u.User.Username()
	params["password"], _ = u.User.Password() // Ignore error for password - ok to be absent

	q := u.Query()
	for k := range q {
		params[k] = q.Get(k)
	}

	params["dbname"] = strings.TrimPrefix(u.Path, "/")

	return params, nil
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	gzipWriter *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gzipWriter.Write(b)
}

func (w *gzipResponseWriter) WriteString(s string) (int, error) {
	return w.gzipWriter.Write([]byte(s))
}

func gzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверка заголовка Content-Encoding
		if !strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Сохранение исходного Content-Type
		originalContentType := c.Writer.Header().Get("Content-Type")
		// Установка заголовка Content-Encoding
		c.Writer.Header().Set("Content-Encoding", "gzip")

		// Создание gzip.Writer
		gw := gzip.NewWriter(c.Writer)
		defer gw.Close()

		// Замена Writer на gzipResponseWriter
		c.Writer = &gzipResponseWriter{
			ResponseWriter: c.Writer,
			gzipWriter:     gw,
		}

		logger, _ := zap.NewDevelopment()
		defer logger.Sync()

		logger.Info("Request processed33ffddd",
			zap.String("method", c.Request.Method),
			zap.String("path", c.ContentType()),
			zap.String("Encoding", c.Request.Header.Get("Accept-Encoding")),
			zap.Int("statusCode", c.Writer.Status()),
		)

		// Вызов следующего обработчика
		c.Next()

		// Восстановление исходного Content-Type
		c.Writer.Header().Set("Content-Type", originalContentType)

		// Закрытие gzipWriter
		if gw, ok := c.Writer.(*gzipResponseWriter); ok {
			gw.gzipWriter.Close()
		}

	}
}

func WithLogging(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()
		start := time.Now()

		h(c)
		//  логи
		duration := time.Since(start)
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", duration), // Получаем статус
			zap.Int("statusCode", c.Writer.Status()),
		)
	}
}

func loadURLsFromFile(fname string) ([]URLData, error) {
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err // Возвращаем nil, чтобы указать, что данные не загружены
	}

	var urls []URLData
	if err := json.Unmarshal(data, &urls); err != nil {
		return nil, err // Возвращаем nil, чтобы указать, что данные не загружены
	}

	return urls, nil // Возвращаем urls, чтобы вернуть загруженные данные
}

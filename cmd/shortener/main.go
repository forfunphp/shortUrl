package main

import (
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"shortUrl/internal/handlers"
	"strings"
	"time"
)

var sugar zap.SugaredLogger

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	router := gin.Default()

	router.Use(gzipMiddleware())
	router.POST("/", WithLogging(handlers.ReduceURL))
	router.GET("/:shortURL", WithLogging(handlers.Redirect))
	router.POST("/api/shorten", WithLogging(handlers.Shorten))

	fmt.Printf("Сервер запущен на %s\n", handlers.Cfg.HTTPAddr)

	colonIndex := strings.Index(handlers.Cfg.HTTPAddr, ":")
	if colonIndex == -1 {
		log.Fatalf("Неверный формат адреса: %s", handlers.Cfg.HTTPAddr)
	}

	port := handlers.Cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
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
		// Проверка заголовка Accept-Encoding
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		logger, _ := zap.NewDevelopment()
		defer logger.Sync()

		logger.Info("Request processed33ff",
			zap.String("method", c.Request.Method),
			zap.String("path", c.ContentType),
			zap.Int("statusCode", c.Writer.Status()),
		)

		// Проверка типа контента
		if c.ContentType() == "application/json" || c.ContentType() == "text/html" {
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

			// Вызов следующего обработчика
			c.Next()

			// Закрытие gzipWriter
			if gw, ok := c.Writer.(*gzipResponseWriter); ok {
				gw.gzipWriter.Close()
			}
		} else {
			c.Next()
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

package main

import (
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

func WithLogging(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()
		start := time.Now()

		h(c)

		// Выводим логи
		duration := time.Since(start)
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", duration), // Получаем статус
		)
	}
}

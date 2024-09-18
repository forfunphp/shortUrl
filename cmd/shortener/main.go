package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortUrl/internal/handlers"
	"strings"
)

func main() {
	router := gin.Default()

	router.POST("/", handlers.ReduceURL)
	router.GET("/:shortURL", handlers.Redirect)

	fmt.Printf("Сервер запущен на %s\n", cfg.HTTPAddr)
	colonIndex := strings.Index(cfg.HTTPAddr, ":")
	port := cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
}

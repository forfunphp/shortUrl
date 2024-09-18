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

	fmt.Printf("Сервер запущен на %s\n", handlers.Cfg.HTTPAddr)

	colonIndex := strings.Index(handlers.Cfg.HTTPAddr, ":")
	if colonIndex == -1 {
		log.Fatalf("Неверный формат адреса: %s", handlers.Cfg.HTTPAddr)
	}

	port := handlers.Cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
}

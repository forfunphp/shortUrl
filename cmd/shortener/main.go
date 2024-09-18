package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shortUrl/config"
	"strings"
)

var cfg = config.NewConfig()

func main() {
	err := cfg.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %v", err)
	}

	router := gin.Default()

	fmt.Printf("Сервер запущен на %s\n", cfg.HTTPAddr)

	colonIndex := strings.Index(cfg.HTTPAddr, ":")
	if colonIndex == -1 {
		log.Fatalf("Неверный формат адреса: %s", cfg.HTTPAddr)
	}

	port := cfg.HTTPAddr[colonIndex:]
	log.Fatal(router.Run(port))
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Сервер запущен")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

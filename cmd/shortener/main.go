package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type URLPair struct {
	Url      string
	ShortURL string
}

var urlMap map[string]URLPair

func reduceURL() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var shortURL string
	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(charset))
		shortURL += string(charset[randomIndex])
	}
	return shortURL
}

func Handler(w http.ResponseWriter, r *http.Request) {

	body := ""
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for k, v := range r.Form {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}

		URL := string(body)
		shortURL := reduceURL()
		urlMap[shortURL] = URLPair{URL, shortURL}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "http://localhost:8080/%s\n", shortURL)

	}

	if r.Method == http.MethodGet {
		u, _ := url.Parse(r.URL.Path)
		parts := strings.Split(u.Path, "/")
		shortURL := strings.Split(parts[1], "favicon.ico")

		urlPair, ok := urlMap[shortURL[0]]
		if !ok {
			http.Error(w, "Нет урла", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", urlPair.Url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

}

func main() {
	urlMap = make(map[string]URLPair)
	http.HandleFunc("/", Handler)
	fmt.Println("Сервер запущен на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

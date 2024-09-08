package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

type URLPair struct {
	URL      *url.URL
	ShortURL string
}

var urlMap = make(map[string]URLPair)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func Handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Не спарсил тело запроса", http.StatusBadRequest)
			return
		}
		URL := string(body)
		parsedURL, err := url.Parse(URL)
		if err != nil {
			http.Error(w, "Не спарсил URL", http.StatusBadRequest)
		}

		shortURL := reduceURL()
		urlMap[shortURL] = URLPair{parsedURL, shortURL}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "http://localhost:8080/%s", shortURL)
	}

	if r.Method == http.MethodGet {
		u, _ := url.Parse(r.URL.String())
		shortURL := strings.Split(u.Path, "/")[1]
		urlPair, ok := urlMap[shortURL]

		if !ok {
			http.Error(w, "Нет урла", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", urlPair.URL.String())
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Сервер запущен на http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

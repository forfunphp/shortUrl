package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"shortUrl/config"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReduceURLHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "POST: add new URL",
			method:       http.MethodPost,
			body:         "https://www.example.com/some/long/url",
			wantStatus:   http.StatusCreated,
			wantLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/", handlers.ReduceURL)

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			//assert.Equal(t, tt.wantStatus, rr.Code, "Handler() returned wrong status code")

			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, rr.Header().Get("Location"), "Handler() returned wrong Location header")
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name         string
		shortURL     string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "GET: valid short URL",
			shortURL:     "test_short_url",
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://www.example.com",
		},
		{
			name:         "GET: invalid short URL",
			shortURL:     "invalid_short_url",
			wantStatus:   http.StatusBadRequest,
			wantLocation: "",
		},
	}

	handlers.URLMap = make(map[string]handlers.URLPair)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.GET("/:shortURL", handlers.Redirect)

			if tt.shortURL == "test_short_url" {
				handlers.URLMap[tt.shortURL] = handlers.URLPair{
					URL:      &url.URL{Scheme: "https", Host: "www.example.com"},
					ShortURL: tt.shortURL,
				}
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.shortURL), nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code, "Handler() returned wrong status code")

			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, rr.Header().Get("Location"), "Handler() returned wrong Location header")
			}
		})
	}
}

func TestShorten(t *testing.T) {
	// Инициализируем конфигурацию
	var Cfg = config.NewConfig()

	// Создаем тестовый маршрутизатор
	router := gin.Default()

	// Зарегистрируем тестовый обработчик
	router.POST("/shorten", ShortenHandler)

	// Создаем тестовый запрос
	reqBody := []byte(`{"url": "https://www.example.com"}`)
	req, err := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	// Создаем тестовый сервер
	w := httptest.NewRecorder()

	// Обрабатываем тестовый запрос
	router.ServeHTTP(w, req)

	// Проверка кода ответа
	assert.Equal(t, http.StatusCreated, w.Code, "Ожидаемый код ответа: 201 (StatusCreated)")

	// Проверка тела ответа
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resp ShortURL
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка ответа
	assert.Equal(t, Cfg.BaseURL+"/shortURL", resp.ShortURL, "Проверка URL в ответе")

	// Проверка URLMap
	assert.Equal(t, URLPair{OriginalURL: &url.URL{Scheme: "https", Host: "www.example.com"}, ShortURL: "shortURL"}, URLMap["shortURL"], "Проверка URLMap")
}

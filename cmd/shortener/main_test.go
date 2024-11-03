package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"shortUrl/internal/handlers"
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

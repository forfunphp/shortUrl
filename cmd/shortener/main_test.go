package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
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
		{
			name:         "GET: redirect to existing URL",
			method:       http.MethodGet,
			body:         "",
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://www.example.com/some/long/url",
		},
		{
			name:         "GET: invalid short URL",
			method:       http.MethodGet,
			body:         "",
			wantStatus:   http.StatusBadRequest,
			wantLocation: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			Handler(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code, "Handler() returned wrong status code")

			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, rr.Header().Get("Location"), "Handler() returned wrong Location header")
			}
		})
	}
}

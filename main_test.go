package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/openimg-go/internal/cache"
)

func TestImageHandler_ServeImage(t *testing.T) {
	handler := &ImageHandler{
		Client: &http.Client{},
		Cache:  cache.New(),
	}

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "valid request",
			url:        "/api/image?url=https://picsum.photos/800/600&w=200&h=200",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing url",
			url:        "/api/image?w=200&h=200",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid dimensions",
			url:        "/api/image?url=https://picsum.photos/800/600&w=5000",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid quality",
			url:        "/api/image?url=https://picsum.photos/800/600&q=101",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid format",
			url:        "/api/image?url=https://picsum.photos/800/600&fmt=gif",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid fit",
			url:        "/api/image?url=https://picsum.photos/800/600&fit=stretch",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			handler.ServeImage(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeImage() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestImageHandler_MethodNotAllowed(t *testing.T) {
	handler := &ImageHandler{
		Client: &http.Client{},
		Cache:  cache.New(),
	}

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/image", nil)
			w := httptest.NewRecorder()
			handler.ServeImage(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("ServeImage() status = %v, want %v", w.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerWithGin(t *testing.T) {
	ServerAddressLong = flag.String("a", "localhost:8080", "HTTP server address")
	ServerAddressShort = flag.String("b", "http://localhost:8080", "Base URL for short links")
	flag.Parse()

	router := gin.Default()

	router.HandleMethodNotAllowed = true
	router.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
	})

	router.GET("/:id", getting)
	router.POST("/", posting)

	t.Run("Create short URL", func(t *testing.T) {
		body := strings.NewReader("https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "text/plain")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		resp := rec.Body.String()
		if !strings.Contains(resp, "http://localhost:8080/") {
			t.Errorf("unexpected response: %s", resp)
		}
	})

	t.Run("Redirect to original URL", func(t *testing.T) {
		id := "test123"
		original := "https://example.com"
		shortToOriginal[id] = original

		req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusTemporaryRedirect {
			t.Errorf("expected %d, got %d", http.StatusTemporaryRedirect, rec.Code)
		}
		loc := rec.Header().Get("Location")
		if loc != original {
			t.Errorf("expected location %s, got %s", original, loc)
		}
	})

	t.Run("Invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

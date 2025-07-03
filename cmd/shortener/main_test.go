package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	// Тестирование POST запроса (создание короткой ссылки)
	t.Run("Create short URL", func(t *testing.T) {
		reqBody := strings.NewReader("https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		req.Header.Set("Content-Type", "text/plain")

		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		shortURL := rec.Body.String()
		if shortURL == "" {
			t.Error("expected non-empty short URL")
		}
	})

	// Тестирование GET запроса (редирект)
	t.Run("Redirect to original URL", func(t *testing.T) {
		// Предварительно создаем тестовую запись
		testID := "test123"
		testURL := "https://example.com"
		shortToOriginal[testID] = testURL

		req := httptest.NewRequest(http.MethodGet, "/"+testID, nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusTemporaryRedirect {
			t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
		}

		location := rec.Header().Get("Location")
		if location != testURL {
			t.Errorf("expected location %s, got %s", testURL, location)
		}
	})

	// Тестирование неверного метода
	t.Run("Invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

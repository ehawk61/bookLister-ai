package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ehawk61/bookLister-ai/internal/middleware"
)

func TestHealthcheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/healthcheck", nil)
	w := httptest.NewRecorder()

	Healthcheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}

func TestHealthcheckWithMiddlewareChain(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthcheck", Healthcheck)
	chain := middleware.CorrelationID(middleware.APIVersion(mux))

	t.Run("valid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/healthcheck", nil)
		req.Header.Set("X-Correlation-ID", "550e8400-e29b-41d4-a716-446655440000")
		w := httptest.NewRecorder()

		chain.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var body map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if body["status"] != "ok" {
			t.Errorf("status = %q, want %q", body["status"], "ok")
		}

		if got := w.Header().Get("X-Correlation-ID"); got != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("X-Correlation-ID = %q, want echoed value", got)
		}
		if got := w.Header().Get("X-API-Version"); got != "1.0" {
			t.Errorf("X-API-Version = %q, want 1.0 (default)", got)
		}
	})

	t.Run("missing correlation ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/healthcheck", nil)
		w := httptest.NewRecorder()

		chain.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

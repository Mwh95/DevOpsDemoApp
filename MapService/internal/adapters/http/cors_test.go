package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	allowed := []string{"http://localhost:5173", "https://example.com"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/public/health/live", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("expected Access-Control-Allow-Origin %q, got %q", "http://localhost:5173", got)
	}
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	allowed := []string{"http://localhost:5173"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/public/health/live", nil)
	req.Header.Set("Origin", "http://evil.example.com")
	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected empty Access-Control-Allow-Origin for disallowed origin, got %q", got)
	}
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	allowed := []string{"http://localhost:5173"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/api/markers", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type")
	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
		t.Errorf("expected 200 or 204 for preflight, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("expected Access-Control-Allow-Origin %q, got %q", "http://localhost:5173", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("expected Access-Control-Allow-Credentials true, got %q", got)
	}
}

func TestCORSMiddleware_PreflightDisallowedOrigin(t *testing.T) {
	allowed := []string{"http://localhost:5173"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/api/markers", nil)
	req.Header.Set("Origin", "http://evil.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected empty Access-Control-Allow-Origin for disallowed origin, got %q", got)
	}
}

func TestCORSMiddleware_MultipleOrigins(t *testing.T) {
	allowed := []string{"http://localhost:5173", "https://app.example.com"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	for _, origin := range allowed {
		req := httptest.NewRequest(http.MethodGet, "/public/health/live", nil)
		req.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()

		srv.Router.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("origin %q: expected Access-Control-Allow-Origin %q, got %q", origin, origin, got)
		}
	}
}

func TestCORSMiddleware_NoOriginHeader(t *testing.T) {
	allowed := []string{"http://localhost:5173"}
	srv, err := NewServer(nil, nil, "", allowed)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/public/health/live", nil)
	rec := httptest.NewRecorder()

	srv.Router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no CORS header without Origin, got %q", got)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for same-origin request, got %d", rec.Code)
	}
}

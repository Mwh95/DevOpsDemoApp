package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_LivenessEndpointsReturnJSON(t *testing.T) {
	srv, err := NewServer(nil, nil, "")
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	for _, path := range []string{"/public/health/live"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		srv.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("%s: expected status 200, got %d", path, rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("%s: expected Content-Type application/json, got %q", path, got)
		}

		var body healthResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("%s: decode response: %v", path, err)
		}
		if body.Status != "UP" {
			t.Fatalf("%s: expected status UP, got %q", path, body.Status)
		}
	}
}

func TestServer_ReadinessEndpointsReturnServiceUnavailableWithoutPool(t *testing.T) {
	srv, err := NewServer(nil, nil, "")
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	for _, path := range []string{"/public/health/ready"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()

		srv.Router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("%s: expected status 503, got %d", path, rec.Code)
		}
		if got := rec.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("%s: expected Content-Type application/json, got %q", path, got)
		}

		var body healthResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("%s: decode response: %v", path, err)
		}
		if body.Status != "DOWN" {
			t.Fatalf("%s: expected status DOWN, got %q", path, body.Status)
		}
	}
}

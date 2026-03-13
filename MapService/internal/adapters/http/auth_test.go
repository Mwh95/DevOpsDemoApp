package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeTokenVerifier implements TokenVerifier for tests.
type fakeTokenVerifier struct {
	subject string
	err    error
}

func (f *fakeTokenVerifier) VerifyAndExtract(ctx context.Context, token string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if token == "valid-token" {
		return f.subject, nil
	}
	return "", errFakeInvalid
}

var errFakeInvalid = errors.New("invalid token")

func TestRequireAuth_NoHeader(t *testing.T) {
	v := &fakeTokenVerifier{subject: "user1"}
	mw := RequireAuth(v)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := mw(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without Authorization, got %d", rec.Code)
	}
}

func TestRequireAuth_InvalidPrefix(t *testing.T) {
	v := &fakeTokenVerifier{subject: "user1"}
	mw := RequireAuth(v)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := mw(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic xyz")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for non-Bearer auth, got %d", rec.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	v := &fakeTokenVerifier{subject: "user1"}
	mw := RequireAuth(v)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	h := mw(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid token, got %d", rec.Code)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	v := &fakeTokenVerifier{subject: "user1"}
	mw := RequireAuth(v)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if UserIDFromContext(r.Context()) != "user1" {
			t.Errorf("expected user1 in context, got %q", UserIDFromContext(r.Context()))
		}
		w.WriteHeader(http.StatusOK)
	})
	h := mw(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 for valid token, got %d", rec.Code)
	}
}

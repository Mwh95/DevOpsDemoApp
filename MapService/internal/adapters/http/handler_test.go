package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/demoapp/map-service/internal/adapters/repository"
	"github.com/demoapp/map-service/internal/domain"
	"github.com/demoapp/map-service/internal/ports"
	"github.com/demoapp/map-service/internal/usecases"
	"github.com/go-chi/chi/v5"
)

// handlerFakeRepo implements ports.MarkerRepository for handler tests.
type handlerFakeRepo struct {
	mu      sync.Mutex
	markers map[string]*domain.Marker
}

func (f *handlerFakeRepo) Create(ctx context.Context, m *domain.Marker) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *m
	f.markers[m.ID] = &copy
	return nil
}

func (f *handlerFakeRepo) GetByID(ctx context.Context, id, userID string) (*domain.Marker, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	m, ok := f.markers[id]
	if !ok || m.UserID != userID {
		return nil, repository.ErrNotFound
	}
	copy := *m
	return &copy, nil
}

func (f *handlerFakeRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Marker, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*domain.Marker
	for _, m := range f.markers {
		if m.UserID == userID {
			copy := *m
			out = append(out, &copy)
		}
	}
	return out, nil
}

func (f *handlerFakeRepo) Update(ctx context.Context, m *domain.Marker) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	existing, ok := f.markers[m.ID]
	if !ok || existing.UserID != m.UserID {
		return repository.ErrNotFound
	}
	copy := *m
	f.markers[m.ID] = &copy
	return nil
}

func (f *handlerFakeRepo) Delete(ctx context.Context, id, userID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	m, ok := f.markers[id]
	if !ok || m.UserID != userID {
		return repository.ErrNotFound
	}
	delete(f.markers, id)
	return nil
}

var _ ports.MarkerRepository = (*handlerFakeRepo)(nil)

func newHandlerFakeRepo() *handlerFakeRepo {
	return &handlerFakeRepo{markers: make(map[string]*domain.Marker)}
}

// testAuthMiddleware injects a fixed user ID into context (for testing without real JWT).
func testAuthMiddleware(userID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TestMarkerHandler_List(t *testing.T) {
	repo := newHandlerFakeRepo()
	uc := usecases.NewMarkerUseCases(repo)
	h := NewMarkerHandler(uc)
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(testAuthMiddleware("user1"))
		r.Get("/markers", h.List)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/markers", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("List: status %d, body %s", rec.Code, rec.Body.String())
	}
	var list []*domain.Marker
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}

	// no user (no auth middleware)
	r2 := chi.NewRouter()
	r2.Get("/api/markers", h.List)
	req2 := httptest.NewRequest(http.MethodGet, "/api/markers", nil)
	rec2 := httptest.NewRecorder()
	r2.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without auth, got %d", rec2.Code)
	}
}

func TestMarkerHandler_Create(t *testing.T) {
	repo := newHandlerFakeRepo()
	uc := usecases.NewMarkerUseCases(repo)
	h := NewMarkerHandler(uc)
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(testAuthMiddleware("user1"))
		r.Post("/markers", h.Create)
	})

	body := `{"latitude":52.52,"longitude":13.405,"label":"Home","note":"My note"}`
	req := httptest.NewRequest(http.MethodPost, "/api/markers", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("Create: status %d, body %s", rec.Code, rec.Body.String())
	}
	var m domain.Marker
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}
	if m.ID == "" || m.UserID != "user1" || m.Label != "Home" || m.Note != "My note" {
		t.Errorf("unexpected marker: %+v", m)
	}

	// invalid body
	reqBad := httptest.NewRequest(http.MethodPost, "/api/markers", bytes.NewReader([]byte("not json")))
	reqBad.Header.Set("Content-Type", "application/json")
	recBad := httptest.NewRecorder()
	r.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", recBad.Code)
	}
}

func TestMarkerHandler_Get(t *testing.T) {
	repo := newHandlerFakeRepo()
	uc := usecases.NewMarkerUseCases(repo)
	h := NewMarkerHandler(uc)
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(testAuthMiddleware("user1"))
		r.Get("/markers/{id}", h.Get)
	})

	// create one
	created, _ := uc.Create(context.Background(), "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 2, Label: "A", Note: ""})

	req := httptest.NewRequest(http.MethodGet, "/api/markers/"+created.ID, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Get: status %d", rec.Code)
	}
	var m domain.Marker
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}
	if m.ID != created.ID {
		t.Errorf("unexpected id: %s", m.ID)
	}

	// not found
	req404 := httptest.NewRequest(http.MethodGet, "/api/markers/nonexistent-id", nil)
	rec404 := httptest.NewRecorder()
	r.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404 for missing marker, got %d", rec404.Code)
	}
}

func TestMarkerHandler_Update(t *testing.T) {
	repo := newHandlerFakeRepo()
	uc := usecases.NewMarkerUseCases(repo)
	h := NewMarkerHandler(uc)
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(testAuthMiddleware("user1"))
		r.Put("/markers/{id}", h.Update)
	})

	created, _ := uc.Create(context.Background(), "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 1, Label: "Old", Note: ""})
	body := `{"label":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/markers/"+created.ID, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Update: status %d, body %s", rec.Code, rec.Body.String())
	}
	var m domain.Marker
	if err := json.NewDecoder(rec.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}
	if m.Label != "Updated" {
		t.Errorf("expected label Updated, got %s", m.Label)
	}

	req404 := httptest.NewRequest(http.MethodPut, "/api/markers/nonexistent", bytes.NewReader([]byte(`{"label":"x"}`)))
	req404.Header.Set("Content-Type", "application/json")
	rec404 := httptest.NewRecorder()
	r.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec404.Code)
	}
}

func TestMarkerHandler_Delete(t *testing.T) {
	repo := newHandlerFakeRepo()
	uc := usecases.NewMarkerUseCases(repo)
	h := NewMarkerHandler(uc)
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(testAuthMiddleware("user1"))
		r.Delete("/markers/{id}", h.Delete)
	})

	created, _ := uc.Create(context.Background(), "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 1, Label: "X", Note: ""})
	req := httptest.NewRequest(http.MethodDelete, "/api/markers/"+created.ID, nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("Delete: status %d", rec.Code)
	}
	list, _ := uc.List(context.Background(), "user1")
	if len(list) != 0 {
		t.Errorf("expected 0 markers after delete, got %d", len(list))
	}

	req404 := httptest.NewRequest(http.MethodDelete, "/api/markers/nonexistent", nil)
	rec404 := httptest.NewRecorder()
	r.ServeHTTP(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec404.Code)
	}
}

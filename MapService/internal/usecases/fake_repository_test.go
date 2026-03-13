package usecases

import (
	"context"
	"sync"

	"github.com/demoapp/map-service/internal/adapters/repository"
	"github.com/demoapp/map-service/internal/domain"
)

// fakeMarkerRepo is an in-memory MarkerRepository for tests. Returns repository.ErrNotFound when not found.
type fakeMarkerRepo struct {
	mu      sync.Mutex
	markers map[string]*domain.Marker
}

func newFakeMarkerRepo() *fakeMarkerRepo {
	return &fakeMarkerRepo{markers: make(map[string]*domain.Marker)}
}

func (f *fakeMarkerRepo) Create(ctx context.Context, m *domain.Marker) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *m
	f.markers[m.ID] = &copy
	return nil
}

func (f *fakeMarkerRepo) GetByID(ctx context.Context, id, userID string) (*domain.Marker, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	m, ok := f.markers[id]
	if !ok || m.UserID != userID {
		return nil, repository.ErrNotFound
	}
	copy := *m
	return &copy, nil
}

func (f *fakeMarkerRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Marker, error) {
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

func (f *fakeMarkerRepo) Update(ctx context.Context, m *domain.Marker) error {
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

func (f *fakeMarkerRepo) Delete(ctx context.Context, id, userID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	m, ok := f.markers[id]
	if !ok || m.UserID != userID {
		return repository.ErrNotFound
	}
	delete(f.markers, id)
	return nil
}

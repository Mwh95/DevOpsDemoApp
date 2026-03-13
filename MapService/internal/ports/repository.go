package ports

import (
	"context"

	"github.com/demoapp/map-service/internal/domain"
)

// MarkerRepository defines persistence for markers.
type MarkerRepository interface {
	Create(ctx context.Context, m *domain.Marker) error
	GetByID(ctx context.Context, id, userID string) (*domain.Marker, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Marker, error)
	Update(ctx context.Context, m *domain.Marker) error
	Delete(ctx context.Context, id, userID string) error
}

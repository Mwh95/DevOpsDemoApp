package usecases

import (
	"context"
	"time"

	"github.com/demoapp/map-service/internal/domain"
	"github.com/demoapp/map-service/internal/ports"
	"github.com/google/uuid"
)

// MarkerUseCases provides application logic for markers.
type MarkerUseCases struct {
	repo ports.MarkerRepository
}

// NewMarkerUseCases returns a new MarkerUseCases.
func NewMarkerUseCases(repo ports.MarkerRepository) *MarkerUseCases {
	return &MarkerUseCases{repo: repo}
}

// Create creates a new marker for the user.
func (uc *MarkerUseCases) Create(ctx context.Context, userID string, in domain.CreateMarkerInput) (*domain.Marker, error) {
	now := time.Now().UTC()
	m := &domain.Marker{
		ID:        uuid.New().String(),
		UserID:    userID,
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
		Label:     in.Label,
		Note:      in.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Get returns a marker by ID for the user.
func (uc *MarkerUseCases) Get(ctx context.Context, id, userID string) (*domain.Marker, error) {
	return uc.repo.GetByID(ctx, id, userID)
}

// List returns all markers for the user.
func (uc *MarkerUseCases) List(ctx context.Context, userID string) ([]*domain.Marker, error) {
	return uc.repo.ListByUserID(ctx, userID)
}

// Update updates a marker's label/note for the user.
func (uc *MarkerUseCases) Update(ctx context.Context, id, userID string, in domain.UpdateMarkerInput) (*domain.Marker, error) {
	m, err := uc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if in.Label != nil {
		m.Label = *in.Label
	}
	if in.Note != nil {
		m.Note = *in.Note
	}
	m.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

// Delete removes a marker for the user.
func (uc *MarkerUseCases) Delete(ctx context.Context, id, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}

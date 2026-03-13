package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/demoapp/map-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestPostgresMarkerRepository_Integration runs against a real Postgres when TEST_DATABASE_URL is set.
func TestPostgresMarkerRepository_Integration(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}

	repo := NewPostgresMarkerRepository(pool)
	userID := "test-user-" + uuid.New().String()[:8]

	m := &domain.Marker{
		ID:        uuid.New().String(),
		UserID:    userID,
		Latitude:  52.52,
		Longitude: 13.405,
		Label:     "Test",
		Note:      "Note",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := repo.Create(ctx, m); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.GetByID(ctx, m.ID, userID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Label != m.Label || got.Note != m.Note {
		t.Errorf("got %+v", got)
	}

	list, err := repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID: %v", err)
	}
	if len(list) < 1 {
		t.Errorf("expected at least 1 marker, got %d", len(list))
	}

	m.Label = "Updated"
	m.Note = "Updated note"
	m.UpdatedAt = time.Now().UTC()
	if err := repo.Update(ctx, m); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, _ = repo.GetByID(ctx, m.ID, userID)
	if got.Label != "Updated" {
		t.Errorf("after update: label %q", got.Label)
	}

	if err := repo.Delete(ctx, m.ID, userID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = repo.GetByID(ctx, m.ID, userID)
	if err != ErrNotFound {
		t.Errorf("after delete: expected ErrNotFound, got %v", err)
	}
}

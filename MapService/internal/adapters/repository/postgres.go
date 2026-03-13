package repository

import (
	"context"
	"errors"

	"github.com/demoapp/map-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("marker not found")

// PostgresMarkerRepository implements ports.MarkerRepository using PostgreSQL.
type PostgresMarkerRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresMarkerRepository returns a new PostgresMarkerRepository.
func NewPostgresMarkerRepository(pool *pgxpool.Pool) *PostgresMarkerRepository {
	return &PostgresMarkerRepository{pool: pool}
}

// Create inserts a marker into mapservice.markers.
func (r *PostgresMarkerRepository) Create(ctx context.Context, m *domain.Marker) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO mapservice.markers (id, user_id, latitude, longitude, label, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, m.ID, m.UserID, m.Latitude, m.Longitude, m.Label, m.Note, m.CreatedAt, m.UpdatedAt)
	return err
}

// GetByID returns a marker by id and user_id.
func (r *PostgresMarkerRepository) GetByID(ctx context.Context, id, userID string) (*domain.Marker, error) {
	var m domain.Marker
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, user_id, latitude, longitude, label, note, created_at, updated_at
		FROM mapservice.markers WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(&m.ID, &m.UserID, &m.Latitude, &m.Longitude, &m.Label, &m.Note, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListByUserID returns all markers for the given user.
func (r *PostgresMarkerRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Marker, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, user_id, latitude, longitude, label, note, created_at, updated_at
		FROM mapservice.markers WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var markers []*domain.Marker
	for rows.Next() {
		var m domain.Marker
		if err := rows.Scan(&m.ID, &m.UserID, &m.Latitude, &m.Longitude, &m.Label, &m.Note, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		markers = append(markers, &m)
	}
	return markers, rows.Err()
}

// Update updates label/note and updated_at for a marker.
func (r *PostgresMarkerRepository) Update(ctx context.Context, m *domain.Marker) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE mapservice.markers SET label = $1, note = $2, updated_at = $3
		WHERE id = $4 AND user_id = $5
	`, m.Label, m.Note, m.UpdatedAt, m.ID, m.UserID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a marker by id and user_id.
func (r *PostgresMarkerRepository) Delete(ctx context.Context, id, userID string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM mapservice.markers WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

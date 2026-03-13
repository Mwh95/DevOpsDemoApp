package domain

import "time"

// Marker represents a user-placed location marker on the map.
type Marker struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Label     string    `json:"label"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMarkerInput holds input for creating a marker.
type CreateMarkerInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Label     string  `json:"label"`
	Note      string  `json:"note"`
}

// UpdateMarkerInput holds input for updating a marker (label/note).
type UpdateMarkerInput struct {
	Label *string `json:"label,omitempty"`
	Note  *string `json:"note,omitempty"`
}

package usecases

import (
	"context"
	"testing"

	"github.com/demoapp/map-service/internal/domain"
)

func TestMarkerUseCases_Create(t *testing.T) {
	repo := newFakeMarkerRepo()
	uc := NewMarkerUseCases(repo)
	ctx := context.Background()

	in := domain.CreateMarkerInput{Latitude: 52.52, Longitude: 13.405, Label: "Home", Note: "My place"}
	m, err := uc.Create(ctx, "user1", in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.UserID != "user1" || m.Latitude != 52.52 || m.Longitude != 13.405 || m.Label != "Home" || m.Note != "My place" {
		t.Errorf("unexpected marker: %+v", m)
	}

	list, _ := uc.List(ctx, "user1")
	if len(list) != 1 {
		t.Errorf("expected 1 marker, got %d", len(list))
	}
}

func TestMarkerUseCases_Get(t *testing.T) {
	repo := newFakeMarkerRepo()
	uc := NewMarkerUseCases(repo)
	ctx := context.Background()

	m, err := uc.Create(ctx, "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 2, Label: "A", Note: ""})
	if err != nil {
		t.Fatal(err)
	}

	got, err := uc.Get(ctx, m.ID, "user1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != m.ID || got.Label != "A" {
		t.Errorf("unexpected: %+v", got)
	}

	_, err = uc.Get(ctx, m.ID, "other-user")
	if err == nil {
		t.Error("expected error for wrong user")
	}
	_, err = uc.Get(ctx, "nonexistent-id", "user1")
	if err == nil {
		t.Error("expected error for missing marker")
	}
}

func TestMarkerUseCases_List(t *testing.T) {
	repo := newFakeMarkerRepo()
	uc := NewMarkerUseCases(repo)
	ctx := context.Background()

	_, _ = uc.Create(ctx, "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 1, Label: "A", Note: ""})
	_, _ = uc.Create(ctx, "user1", domain.CreateMarkerInput{Latitude: 2, Longitude: 2, Label: "B", Note: ""})
	_, _ = uc.Create(ctx, "user2", domain.CreateMarkerInput{Latitude: 3, Longitude: 3, Label: "C", Note: ""})

	list, err := uc.List(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 markers for user1, got %d", len(list))
	}
	list2, _ := uc.List(ctx, "user2")
	if len(list2) != 1 {
		t.Errorf("expected 1 marker for user2, got %d", len(list2))
	}
}

func TestMarkerUseCases_Update(t *testing.T) {
	repo := newFakeMarkerRepo()
	uc := NewMarkerUseCases(repo)
	ctx := context.Background()

	m, _ := uc.Create(ctx, "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 1, Label: "Old", Note: "n1"})
	newLabel := "New"
	updated, err := uc.Update(ctx, m.ID, "user1", domain.UpdateMarkerInput{Label: &newLabel})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Label != "New" || updated.Note != "n1" {
		t.Errorf("unexpected updated: %+v", updated)
	}

	newNote := "n2"
	updated2, _ := uc.Update(ctx, m.ID, "user1", domain.UpdateMarkerInput{Note: &newNote})
	if updated2.Note != "n2" {
		t.Errorf("unexpected note: %+v", updated2)
	}

	_, err = uc.Update(ctx, "nonexistent", "user1", domain.UpdateMarkerInput{Label: &newLabel})
	if err == nil {
		t.Error("expected error for missing marker")
	}
}

func TestMarkerUseCases_Delete(t *testing.T) {
	repo := newFakeMarkerRepo()
	uc := NewMarkerUseCases(repo)
	ctx := context.Background()

	m, _ := uc.Create(ctx, "user1", domain.CreateMarkerInput{Latitude: 1, Longitude: 1, Label: "X", Note: ""})
	err := uc.Delete(ctx, m.ID, "user1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	list, _ := uc.List(ctx, "user1")
	if len(list) != 0 {
		t.Errorf("expected 0 markers after delete, got %d", len(list))
	}

	err = uc.Delete(ctx, "nonexistent", "user1")
	if err == nil {
		t.Error("expected error for missing marker")
	}
}

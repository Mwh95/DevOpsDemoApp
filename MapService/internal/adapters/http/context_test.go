package http

import (
	"context"
	"testing"
)

func TestUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	if got := UserIDFromContext(ctx); got != "" {
		t.Errorf("empty context: got %q", got)
	}
	ctx2 := WithUserID(ctx, "user1")
	if got := UserIDFromContext(ctx2); got != "user1" {
		t.Errorf("with user: got %q", got)
	}
}

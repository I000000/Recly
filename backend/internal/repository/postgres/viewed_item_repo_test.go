//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewedItemRepo_Integration_RecordAndGet(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	viewedRepo := NewViewedItemRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	err = viewedRepo.RecordView(ctx, user.ID, "book", "book-1")
	require.NoError(t, err)
	err = viewedRepo.RecordView(ctx, user.ID, "movie", "movie-1")
	require.NoError(t, err)

	views, err := viewedRepo.GetRecentViews(ctx, user.ID, 10)
	require.NoError(t, err)
	assert.Len(t, views, 2)

	assert.Equal(t, "movie-1", views[0].ItemID)
	assert.Equal(t, "book-1", views[1].ItemID)
}

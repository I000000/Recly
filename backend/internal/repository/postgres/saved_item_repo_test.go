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

func TestSavedItemRepo_Integration_CRUD(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	savedItemRepo := NewSavedItemRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	item, err := savedItemRepo.SaveItem(ctx, user.ID, "book", "book-1")
	require.NoError(t, err)
	assert.NotEmpty(t, item.ID)
	assert.False(t, item.SavedAt.IsZero())

	items, err := savedItemRepo.GetSavedItems(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "book-1", items[0].ItemID)

	err = savedItemRepo.DeleteSavedItem(ctx, item.ID)
	require.NoError(t, err)

	items, err = savedItemRepo.GetSavedItems(ctx, user.ID)
	require.NoError(t, err)
	assert.Empty(t, items)
}

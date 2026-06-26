//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenRepo_Integration_StoreAndGet(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	tokenRepo := NewTokenRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	rt := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: "hashed_token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err = tokenRepo.StoreRefreshToken(ctx, rt)
	require.NoError(t, err)
	assert.NotEmpty(t, rt.ID)
	assert.False(t, rt.CreatedAt.IsZero())

	found, err := tokenRepo.GetRefreshToken(ctx, rt.ID)
	require.NoError(t, err)
	assert.Equal(t, rt.ID, found.ID)
	assert.Equal(t, rt.UserID, found.UserID)
	assert.Equal(t, rt.TokenHash, found.TokenHash)
	assert.WithinDuration(t, rt.ExpiresAt, found.ExpiresAt, time.Second)

	err = tokenRepo.DeleteRefreshToken(ctx, rt.ID)
	require.NoError(t, err)

	_, err = tokenRepo.GetRefreshToken(ctx, rt.ID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

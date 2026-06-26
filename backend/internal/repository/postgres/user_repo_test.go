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

func TestUserRepo_Integration_CreateAndGetByEmail(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "Test User",
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.False(t, user.CreatedAt.IsZero())

	found, err := repo.GetByEmail(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Name, found.Name)
	assert.Equal(t, user.PasswordHash, found.PasswordHash)
	assert.Empty(t, found.AvatarURL)
	assert.False(t, found.OnboardingCompleted)
}

func TestUserRepo_Integration_GetByID(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test2@example.com",
		PasswordHash: "hashed2",
		Name:         "Test User 2",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepo_Integration_UpdateOnboardingCompleted(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test3@example.com",
		PasswordHash: "hashed3",
		Name:         "Test User 3",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.UpdateOnboardingCompleted(ctx, user.ID, true)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.True(t, found.OnboardingCompleted)
}

func TestUserRepo_Integration_UpdateAvatar(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test4@example.com",
		PasswordHash: "hashed4",
		Name:         "Test User 4",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	avatarURL := "http://example.com/avatar.png"
	err = repo.UpdateAvatar(ctx, user.ID, avatarURL)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, avatarURL, found.AvatarURL)
}

func TestUserRepo_Integration_GetByEmail_NotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestUserRepo_Integration_GetByID_NotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepo(pool)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

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

func TestLibraryRepo_Integration_Books(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	libRepo := NewLibraryRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	err = libRepo.AddLikedBook(ctx, user.ID, "book-1")
	require.NoError(t, err)

	err = libRepo.AddLikedBook(ctx, user.ID, "book-2")
	require.NoError(t, err)

	books, err := libRepo.GetLikedBooks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, books, 2)
	assert.ElementsMatch(t, []string{"book-1", "book-2"}, []string{books[0].BookID, books[1].BookID})

	err = libRepo.RemoveLikedBook(ctx, user.ID, "book-1")
	require.NoError(t, err)

	books, err = libRepo.GetLikedBooks(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "book-2", books[0].BookID)
}

func TestLibraryRepo_Integration_Movies(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	libRepo := NewLibraryRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	err = libRepo.AddLikedMovie(ctx, user.ID, "movie-1")
	require.NoError(t, err)
	err = libRepo.AddLikedMovie(ctx, user.ID, "movie-2")
	require.NoError(t, err)

	movies, err := libRepo.GetLikedMovies(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, movies, 2)

	err = libRepo.RemoveLikedMovie(ctx, user.ID, "movie-1")
	require.NoError(t, err)

	movies, err = libRepo.GetLikedMovies(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, movies, 1)
	assert.Equal(t, "movie-2", movies[0].MovieID)
}

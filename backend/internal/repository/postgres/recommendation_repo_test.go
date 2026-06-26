//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecommendationRepo_Integration_History(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := NewUserRepo(pool)
	recRepo := NewRecommendationRepo(pool)
	ctx := context.Background()

	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test User",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	taskID := uuid.New().String()
	entry := &domain.RecommendationHistory{
		UserID:      user.ID,
		TaskID:      taskID,
		SelectedIDs: []string{"book_1", "movie_1"},
		Direction:   "book_to_movie",
		Weights:     `{"text":0.4,"genre":0.3,"image":0.3}`,
	}
	err = recRepo.SaveHistory(ctx, entry)
	require.NoError(t, err)
	assert.NotEmpty(t, entry.ID)
	assert.False(t, entry.CreatedAt.IsZero())

	history, err := recRepo.GetHistory(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, history, 1)
	assert.Equal(t, taskID, history[0].TaskID)

	found, err := recRepo.GetHistoryByTaskID(ctx, taskID)
	require.NoError(t, err)
	assert.Equal(t, entry.ID, found.ID)

	resultJSON := `["movie_1", "movie_2"]`
	err = recRepo.UpdateResult(ctx, taskID, resultJSON)
	require.NoError(t, err)

	updated, err := recRepo.GetHistoryByTaskID(ctx, taskID)
	require.NoError(t, err)
	assert.Equal(t, resultJSON, updated.Result)
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSourceSelector_Select_ExplicitIDs(t *testing.T) {
	libRepo := mocks.NewLibraryRepository(t)
	libSvc := NewLibraryService(libRepo)
	selector := NewSourceSelector(libSvc)

	ids := []string{"book_1", "movie_1"}
	selected, weights, err := selector.Select(context.Background(), "user-1", ids)

	assert.NoError(t, err)
	assert.Equal(t, ids, selected)
	assert.Len(t, weights, 2)
	assert.Equal(t, 1.0, weights["book_1"])
	assert.Equal(t, 1.0, weights["movie_1"])
}

func TestSourceSelector_Select_FromLibrary(t *testing.T) {
	libRepo := mocks.NewLibraryRepository(t)
	libSvc := NewLibraryService(libRepo)

	now := time.Now()
	books := []domain.LikedBook{
		{UserID: "user-1", BookID: "1", LikedAt: now.Add(-time.Hour)},
		{UserID: "user-1", BookID: "2", LikedAt: now.Add(-24 * time.Hour)},
	}
	movies := []domain.LikedMovie{
		{UserID: "user-1", MovieID: "1", LikedAt: now.Add(-2 * time.Hour)},
	}

	libRepo.On("GetLikedBooks", mock.Anything, "user-1").Return(books, nil)
	libRepo.On("GetLikedMovies", mock.Anything, "user-1").Return(movies, nil)

	selector := NewSourceSelector(libSvc)

	selected, weights, err := selector.Select(context.Background(), "user-1", nil)

	assert.NoError(t, err)
	assert.Len(t, selected, 3)
	assert.Contains(t, selected, "book_1")
	assert.Contains(t, selected, "book_2")
	assert.Contains(t, selected, "movie_1")

	assert.True(t, weights["book_1"] > weights["book_2"], "book_1 should have greater weight than book_2")
	assert.True(t, weights["movie_1"] > weights["book_2"], "movie_1 should have greater weight than book_2")
}

func TestSourceSelector_Select_NoLibraryItems(t *testing.T) {
	libRepo := mocks.NewLibraryRepository(t)
	libSvc := NewLibraryService(libRepo)

	libRepo.On("GetLikedBooks", mock.Anything, "user-1").Return([]domain.LikedBook{}, nil)
	libRepo.On("GetLikedMovies", mock.Anything, "user-1").Return([]domain.LikedMovie{}, nil)

	selector := NewSourceSelector(libSvc)

	selected, weights, err := selector.Select(context.Background(), "user-1", nil)

	assert.Error(t, err)
	assert.Nil(t, selected)
	assert.Nil(t, weights)
	libRepo.AssertExpectations(t)
}

func TestSourceSelector_Select_LibraryError(t *testing.T) {
	libRepo := mocks.NewLibraryRepository(t)
	libSvc := NewLibraryService(libRepo)

	libRepo.On("GetLikedBooks", mock.Anything, "user-1").Return(nil, assert.AnError)

	selector := NewSourceSelector(libSvc)

	selected, weights, err := selector.Select(context.Background(), "user-1", nil)

	assert.Error(t, err)
	assert.Nil(t, selected)
	assert.Nil(t, weights)
	libRepo.AssertExpectations(t)
}

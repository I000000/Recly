package service

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLibraryService_AddBook_Success(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("AddLikedBook", mock.Anything, "user-1", "book-1").Return(nil)

	svc := NewLibraryService(repo)
	err := svc.AddBook(context.Background(), "user-1", "book-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLibraryService_RemoveBook_Success(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("RemoveLikedBook", mock.Anything, "user-1", "book-1").Return(nil)

	svc := NewLibraryService(repo)
	err := svc.RemoveBook(context.Background(), "user-1", "book-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLibraryService_GetBooks_Success(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	expected := []domain.LikedBook{
		{UserID: "user-1", BookID: "book-1"},
		{UserID: "user-1", BookID: "book-2"},
	}
	repo.On("GetLikedBooks", mock.Anything, "user-1").Return(expected, nil)

	svc := NewLibraryService(repo)
	books, err := svc.GetBooks(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, "book-1", books[0].BookID)
	repo.AssertExpectations(t)
}

// Аналогично для фильмов (AddMovie, RemoveMovie, GetMovies)
func TestLibraryService_AddMovie_Success(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("AddLikedMovie", mock.Anything, "user-1", "movie-1").Return(nil)

	svc := NewLibraryService(repo)
	err := svc.AddMovie(context.Background(), "user-1", "movie-1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLibraryService_AddBook_Error(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("AddLikedBook", mock.Anything, "user-1", "book-1").Return(assert.AnError)

	svc := NewLibraryService(repo)
	err := svc.AddBook(context.Background(), "user-1", "book-1")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestLibraryService_GetBooks_Error(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("GetLikedBooks", mock.Anything, "user-1").Return(nil, assert.AnError)

	svc := NewLibraryService(repo)
	books, err := svc.GetBooks(context.Background(), "user-1")
	assert.Error(t, err)
	assert.Nil(t, books)
	repo.AssertExpectations(t)
}

func TestLibraryService_AddMovie_Error(t *testing.T) {
	repo := mocks.NewLibraryRepository(t)
	repo.On("AddLikedMovie", mock.Anything, "user-1", "movie-1").Return(assert.AnError)

	svc := NewLibraryService(repo)
	err := svc.AddMovie(context.Background(), "user-1", "movie-1")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

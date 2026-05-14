package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/service"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLibraryService_AddBook(t *testing.T) {
	repo := &mocks.LibraryRepository{}
	repo.On("AddLikedBook", context.Background(), "user1", "book1").Return(nil)
	svc := service.NewLibraryService(repo)
	err := svc.AddBook(context.Background(), "user1", "book1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLibraryService_GetBooks(t *testing.T) {
	repo := &mocks.LibraryRepository{}
	expected := []domain.LikedBook{
		{UserID: "user1", BookID: "book1", LikedAt: time.Now()},
	}
	repo.On("GetLikedBooks", context.Background(), "user1").Return(expected, nil)
	svc := service.NewLibraryService(repo)
	books, err := svc.GetBooks(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "book1", books[0].BookID)
	repo.AssertExpectations(t)
}

package service

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSavedItemService_SaveItem_Success(t *testing.T) {
	repo := mocks.NewSavedItemRepository(t)
	expected := &domain.SavedItem{ID: "saved-1", UserID: "user-1", ItemType: "book", ItemID: "book-1"}
	repo.On("SaveItem", mock.Anything, "user-1", "book", "book-1").Return(expected, nil)

	svc := NewSavedItemService(repo)
	item, err := svc.SaveItem(context.Background(), "user-1", "book", "book-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, item)
	repo.AssertExpectations(t)
}

func TestSavedItemService_DeleteSavedItem_Success(t *testing.T) {
	repo := mocks.NewSavedItemRepository(t)
	repo.On("DeleteSavedItem", mock.Anything, "saved-1").Return(nil)

	svc := NewSavedItemService(repo)
	err := svc.DeleteSavedItem(context.Background(), "saved-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestSavedItemService_SaveItem_Error(t *testing.T) {
	repo := mocks.NewSavedItemRepository(t)
	repo.On("SaveItem", mock.Anything, "user-1", "book", "book-1").Return(nil, assert.AnError)

	svc := NewSavedItemService(repo)
	item, err := svc.SaveItem(context.Background(), "user-1", "book", "book-1")
	assert.Error(t, err)
	assert.Nil(t, item)
	repo.AssertExpectations(t)
}

func TestSavedItemService_DeleteSavedItem_Error(t *testing.T) {
	repo := mocks.NewSavedItemRepository(t)
	repo.On("DeleteSavedItem", mock.Anything, "saved-1").Return(assert.AnError)

	svc := NewSavedItemService(repo)
	err := svc.DeleteSavedItem(context.Background(), "saved-1")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

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

func TestViewedItemService_RecordView_Success(t *testing.T) {
	repo := mocks.NewViewedItemRepository(t)
	repo.On("RecordView", mock.Anything, "user-1", "book", "book-1").Return(nil)

	svc := NewViewedItemService(repo)
	err := svc.RecordView(context.Background(), "user-1", "book", "book-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestViewedItemService_GetRecentViews_Success(t *testing.T) {
	repo := mocks.NewViewedItemRepository(t)
	expected := []domain.ViewedItem{
		{ID: "view-1", UserID: "user-1", ItemType: "book", ItemID: "book-1", ViewedAt: time.Now()},
	}
	repo.On("GetRecentViews", mock.Anything, "user-1", 20).Return(expected, nil)

	svc := NewViewedItemService(repo)
	views, err := svc.GetRecentViews(context.Background(), "user-1", 20)

	assert.NoError(t, err)
	assert.Len(t, views, 1)
	repo.AssertExpectations(t)
}

func TestViewedItemService_RecordView_Error(t *testing.T) {
	repo := mocks.NewViewedItemRepository(t)
	repo.On("RecordView", mock.Anything, "user-1", "book", "book-1").Return(assert.AnError)

	svc := NewViewedItemService(repo)
	err := svc.RecordView(context.Background(), "user-1", "book", "book-1")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestViewedItemService_GetRecentViews_Error(t *testing.T) {
	repo := mocks.NewViewedItemRepository(t)
	repo.On("GetRecentViews", mock.Anything, "user-1", 20).Return(nil, assert.AnError)

	svc := NewViewedItemService(repo)
	views, err := svc.GetRecentViews(context.Background(), "user-1", 20)
	assert.Error(t, err)
	assert.Nil(t, views)
	repo.AssertExpectations(t)
}

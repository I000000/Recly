package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type ViewedItemService struct {
	repo domain.ViewedItemRepository
}

func NewViewedItemService(repo domain.ViewedItemRepository) *ViewedItemService {
	return &ViewedItemService{repo: repo}
}

func (s *ViewedItemService) RecordView(ctx context.Context, userID, itemType, itemID string) error {
	return s.repo.RecordView(ctx, userID, itemType, itemID)
}

func (s *ViewedItemService) GetRecentViews(ctx context.Context, userID string, limit int) ([]domain.ViewedItem, error) {
	return s.repo.GetRecentViews(ctx, userID, limit)
}

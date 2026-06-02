package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type SavedItemService struct {
	repo domain.SavedItemRepository
}

func NewSavedItemService(repo domain.SavedItemRepository) *SavedItemService {
	return &SavedItemService{repo: repo}
}

func (s *SavedItemService) SaveItem(ctx context.Context, userID, itemType, itemID string) (*domain.SavedItem, error) {
	return s.repo.SaveItem(ctx, userID, itemType, itemID)
}

func (s *SavedItemService) DeleteSavedItem(ctx context.Context, id string) error {
	return s.repo.DeleteSavedItem(ctx, id)
}

func (s *SavedItemService) GetSavedItems(ctx context.Context, userID string) ([]domain.SavedItem, error) {
	return s.repo.GetSavedItems(ctx, userID)
}

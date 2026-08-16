package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type SavedItemService struct {
	repo domain.SavedItemRepository
}

func NewSavedItemService(repo domain.SavedItemRepository) *SavedItemService {
	return &SavedItemService{repo: repo}
}

func (s *SavedItemService) SaveItem(ctx context.Context, userID, itemType, itemID string) (*domain.SavedItem, error) {
	log := logger.FromContext(ctx)
	log.Info("saving item", zap.String("user_id", userID), zap.String("item_type", itemType), zap.String("item_id", itemID))
	item, err := s.repo.SaveItem(ctx, userID, itemType, itemID)
	if err != nil {
		log.Error("failed to save item", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("item saved", zap.String("user_id", userID), zap.String("saved_id", item.ID))
	return item, nil
}

func (s *SavedItemService) DeleteSavedItem(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Info("deleting saved item", zap.String("id", id))
	if err := s.repo.DeleteSavedItem(ctx, id); err != nil {
		log.Error("failed to delete saved item", zap.Error(err), zap.String("id", id))
		return err
	}
	log.Debug("saved item deleted", zap.String("id", id))
	return nil
}

func (s *SavedItemService) GetSavedItems(ctx context.Context, userID string) ([]domain.SavedItem, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching saved items", zap.String("user_id", userID))
	items, err := s.repo.GetSavedItems(ctx, userID)
	if err != nil {
		log.Error("failed to get saved items", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("saved items fetched", zap.String("user_id", userID), zap.Int("count", len(items)))
	return items, nil
}

package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type ViewedItemService struct {
	repo domain.ViewedItemRepository
}

func NewViewedItemService(repo domain.ViewedItemRepository) *ViewedItemService {
	return &ViewedItemService{repo: repo}
}

func (s *ViewedItemService) RecordView(ctx context.Context, userID, itemType, itemID string) error {
	log := logger.FromContext(ctx)
	log.Info("recording view", zap.String("user_id", userID), zap.String("item_type", itemType), zap.String("item_id", itemID))
	if err := s.repo.RecordView(ctx, userID, itemType, itemID); err != nil {
		log.Error("failed to record view", zap.Error(err), zap.String("user_id", userID))
		return err
	}
	log.Debug("view recorded", zap.String("user_id", userID), zap.String("item_id", itemID))
	return nil
}

func (s *ViewedItemService) GetRecentViews(ctx context.Context, userID string, limit int) ([]domain.ViewedItem, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching recent views", zap.String("user_id", userID), zap.Int("limit", limit))
	views, err := s.repo.GetRecentViews(ctx, userID, limit)
	if err != nil {
		log.Error("failed to get recent views", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("recent views fetched", zap.String("user_id", userID), zap.Int("count", len(views)))
	return views, nil
}

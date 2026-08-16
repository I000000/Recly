package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type ViewedItemRepo struct {
	pool Querier
}

func NewViewedItemRepo(pool Querier) domain.ViewedItemRepository {
	return &ViewedItemRepo{pool: pool}
}

func (r *ViewedItemRepo) RecordView(ctx context.Context, userID, itemType, itemID string) error {
	log := logger.FromContext(ctx)
	log.Debug("recording view", zap.String("user_id", userID), zap.String("item_type", itemType), zap.String("item_id", itemID))
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_viewed_items (user_id, item_type, item_id, viewed_at)
         VALUES ($1, $2, $3, NOW())
         ON CONFLICT (user_id, item_type, item_id) DO UPDATE SET viewed_at = NOW()`,
		userID, itemType, itemID,
	)
	if err != nil {
		log.Error("failed to record view", zap.Error(err), zap.String("user_id", userID))
	}
	return err
}

func (r *ViewedItemRepo) GetRecentViews(ctx context.Context, userID string, limit int) ([]domain.ViewedItem, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting recent views", zap.String("user_id", userID), zap.Int("limit", limit))
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, item_type, item_id, viewed_at
         FROM user_viewed_items
         WHERE user_id = $1
         ORDER BY viewed_at DESC
         LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		log.Error("failed to query recent views", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var items []domain.ViewedItem
	for rows.Next() {
		var v domain.ViewedItem
		if err := rows.Scan(&v.ID, &v.UserID, &v.ItemType, &v.ItemID, &v.ViewedAt); err != nil {
			log.Error("failed to scan viewed item", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("recent views retrieved", zap.String("user_id", userID), zap.Int("count", len(items)))
	return items, nil
}

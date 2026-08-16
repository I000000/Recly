package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type SavedItemRepo struct {
	pool Querier
}

func NewSavedItemRepo(pool Querier) *SavedItemRepo {
	return &SavedItemRepo{pool: pool}
}

func (r *SavedItemRepo) SaveItem(ctx context.Context, userID, itemType, itemID string) (*domain.SavedItem, error) {
	log := logger.FromContext(ctx)
	log.Debug("saving item", zap.String("user_id", userID), zap.String("item_type", itemType), zap.String("item_id", itemID))
	item := &domain.SavedItem{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO user_saved_items (user_id, item_type, item_id)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id, item_type, item_id) DO NOTHING
         RETURNING id, user_id, item_type, item_id, saved_at`,
		userID, itemType, itemID,
	).Scan(&item.ID, &item.UserID, &item.ItemType, &item.ItemID, &item.SavedAt)
	if err != nil {
		log.Error("failed to save item", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("item saved", zap.String("user_id", userID), zap.String("saved_id", item.ID))
	return item, nil
}

func (r *SavedItemRepo) DeleteSavedItem(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Debug("deleting saved item", zap.String("id", id))
	_, err := r.pool.Exec(ctx, `DELETE FROM user_saved_items WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete saved item", zap.Error(err), zap.String("id", id))
	}
	return err
}

func (r *SavedItemRepo) GetSavedItems(ctx context.Context, userID string) ([]domain.SavedItem, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting saved items", zap.String("user_id", userID))
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, item_type, item_id, saved_at
         FROM user_saved_items WHERE user_id = $1 ORDER BY saved_at DESC`, userID)
	if err != nil {
		log.Error("failed to query saved items", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var items []domain.SavedItem
	for rows.Next() {
		var item domain.SavedItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ItemType, &item.ItemID, &item.SavedAt); err != nil {
			log.Error("failed to scan saved item", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("saved items retrieved", zap.String("user_id", userID), zap.Int("count", len(items)))
	return items, nil
}

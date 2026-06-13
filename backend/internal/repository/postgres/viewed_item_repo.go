package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type ViewedItemRepo struct {
	pool Querier
}

func NewViewedItemRepo(pool Querier) domain.ViewedItemRepository {
	return &ViewedItemRepo{pool: pool}
}

func (r *ViewedItemRepo) RecordView(ctx context.Context, userID, itemType, itemID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_viewed_items (user_id, item_type, item_id, viewed_at)
         VALUES ($1, $2, $3, NOW())
         ON CONFLICT (user_id, item_type, item_id) DO UPDATE SET viewed_at = NOW()`,
		userID, itemType, itemID,
	)
	return err
}

func (r *ViewedItemRepo) GetRecentViews(ctx context.Context, userID string, limit int) ([]domain.ViewedItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, item_type, item_id, viewed_at
         FROM user_viewed_items
         WHERE user_id = $1
         ORDER BY viewed_at DESC
         LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ViewedItem
	for rows.Next() {
		var v domain.ViewedItem
		err := rows.Scan(&v.ID, &v.UserID, &v.ItemType, &v.ItemID, &v.ViewedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return items, rows.Err()
}

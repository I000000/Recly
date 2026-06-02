package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type SavedItemRepo struct {
	pool Querier
}

func NewSavedItemRepo(pool Querier) *SavedItemRepo {
	return &SavedItemRepo{pool: pool}
}

func (r *SavedItemRepo) SaveItem(ctx context.Context, userID, itemType, itemID string) (*domain.SavedItem, error) {
	item := &domain.SavedItem{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO user_saved_items (user_id, item_type, item_id)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id, item_type, item_id) DO NOTHING
         RETURNING id, user_id, item_type, item_id, saved_at`,
		userID, itemType, itemID,
	).Scan(&item.ID, &item.UserID, &item.ItemType, &item.ItemID, &item.SavedAt)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *SavedItemRepo) DeleteSavedItem(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM user_saved_items WHERE id = $1`, id)
	return err
}

func (r *SavedItemRepo) GetSavedItems(ctx context.Context, userID string) ([]domain.SavedItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, item_type, item_id, saved_at
         FROM user_saved_items WHERE user_id = $1 ORDER BY saved_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.SavedItem
	for rows.Next() {
		var item domain.SavedItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ItemType, &item.ItemID, &item.SavedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

package domain

import (
	"context"
	"time"
)

type SavedItem struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	ItemType string    `json:"item_type"`
	ItemID   string    `json:"item_id"`
	SavedAt  time.Time `json:"saved_at"`
}

type SavedItemRepository interface {
	SaveItem(ctx context.Context, userID, itemType, itemID string) (*SavedItem, error)
	DeleteSavedItem(ctx context.Context, id string) error
	GetSavedItems(ctx context.Context, userID string) ([]SavedItem, error)
}

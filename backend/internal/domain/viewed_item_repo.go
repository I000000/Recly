package domain

import "context"

type ViewedItemRepository interface {
	RecordView(ctx context.Context, userID, itemType, itemID string) error
	GetRecentViews(ctx context.Context, userID string, limit int) ([]ViewedItem, error)
}

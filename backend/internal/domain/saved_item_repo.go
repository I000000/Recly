//go:generate mockery --name SavedItemRepository --output ../../mocks --outpkg mocks --case underscore
package domain

import "context"

type SavedItemRepository interface {
	SaveItem(ctx context.Context, userID, itemType, itemID string) (*SavedItem, error)
	DeleteSavedItem(ctx context.Context, id string) error
	GetSavedItems(ctx context.Context, userID string) ([]SavedItem, error)
}

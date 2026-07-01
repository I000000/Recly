//go:generate mockery --name SavedItemService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type SavedItemService interface {
	SaveItem(ctx context.Context, userID, itemType, itemID string) (*domain.SavedItem, error)
	DeleteSavedItem(ctx context.Context, id string) error
	GetSavedItems(ctx context.Context, userID string) ([]domain.SavedItem, error)
}

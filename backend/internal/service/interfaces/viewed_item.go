//go:generate mockery --name ViewedItemService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type ViewedItemService interface {
	RecordView(ctx context.Context, userID, itemType, itemID string) error
	GetRecentViews(ctx context.Context, userID string, limit int) ([]domain.ViewedItem, error)
}

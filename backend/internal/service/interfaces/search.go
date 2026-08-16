//go:generate mockery --name SearchService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type SearchService interface {
	Search(ctx context.Context, query string) ([]domain.ItemDetail, error)
	SearchWithFilters(ctx context.Context, query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error)
	GetItems(ctx context.Context, ids []string, itemType string) ([]domain.ItemDetail, error)
	GetGenres(ctx context.Context, itemType string) ([]string, error)
}

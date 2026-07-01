//go:generate mockery --name SearchService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"github.com/I000000/recly/internal/domain"
)

type SearchService interface {
	Search(query string) ([]domain.ItemDetail, error)
	SearchWithFilters(query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error)
	GetItems(ids []string, itemType string) ([]domain.ItemDetail, error)
	GetGenres(itemType string) ([]string, error)
}

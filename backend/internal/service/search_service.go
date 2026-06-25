package service

import (
	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/meili"
)

type SearchService struct {
	meiliClient *meili.Client
}

func NewSearchService(meiliClient *meili.Client) *SearchService {
	return &SearchService{meiliClient: meiliClient}
}

func (s *SearchService) Search(query string) ([]domain.ItemDetail, error) {
	return s.meiliClient.Search(query)
}

func (s *SearchService) SearchWithFilters(query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error) {
	return s.meiliClient.SearchWithFilters(query, itemType, genre, sort, limit, offset)
}

func (s *SearchService) GetItems(ids []string, itemType string) ([]domain.ItemDetail, error) {
	return s.meiliClient.GetItems(ids, itemType)
}

func (s *SearchService) GetGenres(itemType string) ([]string, error) {
	return s.meiliClient.GetGenres(itemType)
}

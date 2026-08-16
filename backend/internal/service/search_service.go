package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/meili"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type SearchService struct {
	meiliClient meili.Client
}

func NewSearchService(meiliClient meili.Client) *SearchService {
	return &SearchService{meiliClient: meiliClient}
}

func (s *SearchService) Search(ctx context.Context, query string) ([]domain.ItemDetail, error) {
	log := logger.FromContext(ctx)
	log.Debug("performing search", zap.String("query", query))
	results, err := s.meiliClient.Search(ctx, query)
	if err != nil {
		log.Error("search failed", zap.Error(err), zap.String("query", query))
		return nil, err
	}
	log.Debug("search completed", zap.String("query", query), zap.Int("count", len(results)))
	return results, nil
}

func (s *SearchService) SearchWithFilters(ctx context.Context, query, itemType, genre, sort string, limit, offset int) ([]domain.ItemDetail, error) {
	log := logger.FromContext(ctx)
	log.Debug("searching with filters", zap.String("query", query), zap.String("type", itemType), zap.String("genre", genre))
	results, err := s.meiliClient.SearchWithFilters(ctx, query, itemType, genre, sort, limit, offset)
	if err != nil {
		log.Error("search with filters failed", zap.Error(err), zap.String("query", query))
		return nil, err
	}
	log.Debug("search with filters completed", zap.String("query", query), zap.Int("count", len(results)))
	return results, nil
}

func (s *SearchService) GetItems(ctx context.Context, ids []string, itemType string) ([]domain.ItemDetail, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting items", zap.Int("count", len(ids)), zap.String("type", itemType))
	items, err := s.meiliClient.GetItems(ctx, ids, itemType)
	if err != nil {
		log.Error("failed to get items", zap.Error(err), zap.Strings("ids", ids))
		return nil, err
	}
	log.Debug("items fetched", zap.Int("count", len(items)))
	return items, nil
}

func (s *SearchService) GetGenres(ctx context.Context, itemType string) ([]string, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting genres", zap.String("type", itemType))
	genres, err := s.meiliClient.GetGenres(ctx, itemType)
	if err != nil {
		log.Error("failed to get genres", zap.Error(err), zap.String("type", itemType))
		return nil, err
	}
	log.Debug("genres fetched", zap.Int("count", len(genres)))
	return genres, nil
}

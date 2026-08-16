package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type SourceSelector struct {
	libSvc *LibraryService
}

func NewSourceSelector(libSvc *LibraryService) *SourceSelector {
	return &SourceSelector{libSvc: libSvc}
}

func (s *SourceSelector) Select(ctx context.Context, userID string, ids []string) ([]string, map[string]float64, error) {
	log := logger.FromContext(ctx)
	if len(ids) > 0 {
		log.Debug("using explicit IDs", zap.String("user_id", userID), zap.Int("count", len(ids)))
		weights := make(map[string]float64, len(ids))
		for _, id := range ids {
			weights[id] = 1.0
		}
		return ids, weights, nil
	}

	log.Debug("selecting sources from library", zap.String("user_id", userID))
	books, err := s.libSvc.GetBooks(ctx, userID)
	if err != nil {
		log.Error("failed to get liked books", zap.Error(err), zap.String("user_id", userID))
		return nil, nil, err
	}
	movies, err := s.libSvc.GetMovies(ctx, userID)
	if err != nil {
		log.Error("failed to get liked movies", zap.Error(err), zap.String("user_id", userID))
		return nil, nil, err
	}

	if len(books) == 0 && len(movies) == 0 {
		log.Warn("no liked items found", zap.String("user_id", userID))
		return nil, nil, errors.New("no liked items to recommend from")
	}

	var selectedIDs []string
	weights := make(map[string]float64)
	tau := 30 * 24 * time.Hour
	now := time.Now()

	for _, b := range books {
		key := "book_" + b.BookID
		selectedIDs = append(selectedIDs, key)
		age := now.Sub(b.LikedAt)
		weights[key] = math.Exp(-age.Seconds() / tau.Seconds())
	}
	for _, m := range movies {
		key := "movie_" + m.MovieID
		selectedIDs = append(selectedIDs, key)
		age := now.Sub(m.LikedAt)
		weights[key] = math.Exp(-age.Seconds() / tau.Seconds())
	}
	log.Debug("sources selected", zap.String("user_id", userID), zap.Int("count", len(selectedIDs)))
	return selectedIDs, weights, nil
}

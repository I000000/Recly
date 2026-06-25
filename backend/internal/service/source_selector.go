package service

import (
	"context"
	"errors"
	"math"
	"time"
)

type SourceSelector struct {
	libSvc *LibraryService
}

func NewSourceSelector(libSvc *LibraryService) *SourceSelector {
	return &SourceSelector{libSvc: libSvc}
}

func (s *SourceSelector) Select(ctx context.Context, userID string, ids []string) ([]string, map[string]float64, error) {
	if len(ids) > 0 {
		weights := make(map[string]float64, len(ids))
		for _, id := range ids {
			weights[id] = 1.0
		}
		return ids, weights, nil
	}

	books, err := s.libSvc.GetBooks(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	movies, err := s.libSvc.GetMovies(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(books) == 0 && len(movies) == 0 {
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

	return selectedIDs, weights, nil
}

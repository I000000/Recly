package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/rabbitmq"
	"github.com/I000000/recly/internal/redis"
	"github.com/google/uuid"
)

type RecommendationService struct {
	repo      domain.RecommendationRepository
	publisher rabbitmq.Publisher
	cache     redis.Cache
	libSvc    *LibraryService
}

func NewRecommendationService(
	repo domain.RecommendationRepository,
	pub rabbitmq.Publisher,
	cache redis.Cache,
	libSvc *LibraryService,
) *RecommendationService {
	return &RecommendationService{repo: repo, publisher: pub, cache: cache, libSvc: libSvc}
}

func (s *RecommendationService) Request(ctx context.Context, userID string, selectedIDs []string, weights map[string]float64, excludeIDs []string, direction string, contextual bool) (string, error) {

	var selectedWeights map[string]float64

	if len(selectedIDs) == 0 {
		books, _ := s.libSvc.GetBooks(ctx, userID)
		movies, _ := s.libSvc.GetMovies(ctx, userID)

		tau := 30 * 24 * time.Hour
		now := time.Now()
		selectedWeights = make(map[string]float64)

		for _, b := range books {
			key := "book_" + b.BookID
			selectedIDs = append(selectedIDs, key)
			age := now.Sub(b.LikedAt)
			selectedWeights[key] = math.Exp(-age.Seconds() / tau.Seconds())
		}
		for _, m := range movies {
			key := "movie_" + m.MovieID
			selectedIDs = append(selectedIDs, key)
			age := now.Sub(m.LikedAt)
			selectedWeights[key] = math.Exp(-age.Seconds() / tau.Seconds())
		}
		if len(selectedIDs) == 0 {
			return "", errors.New("no liked items to recommend from")
		}
	}

	taskID := uuid.New().String()
	msg := rabbitmq.TaskMessage{
		TaskID:          taskID,
		UserID:          userID,
		SelectedIDs:     selectedIDs,
		SelectedWeights: selectedWeights,
		ExcludeIDs:      excludeIDs,
		Direction:       direction,
		Weights:         weights,
		Contextual:      contextual,
	}
	if err := s.publisher.PublishRecommendationTask(ctx, msg); err != nil {
		return "", err
	}
	if err := s.cache.SetResult(ctx, taskID, redis.RecommendationResult{
		Status:     "pending",
		CreatedAt:  time.Now().Unix(),
		Contextual: contextual, // ← добавить
	}, 30*time.Minute); err != nil {
		// не фатально
	}

	wJSON, _ := json.Marshal(weights)
	entry := &domain.RecommendationHistory{
		UserID:      userID,
		TaskID:      taskID,
		SelectedIDs: selectedIDs,
		Direction:   direction,
		Weights:     string(wJSON),
	}
	if err := s.repo.SaveHistory(ctx, entry); err != nil {
		return "", err
	}
	return taskID, nil
}

func (s *RecommendationService) GetResult(ctx context.Context, taskID string) (*redis.RecommendationResult, error) {
	result, err := s.cache.GetResult(ctx, taskID)
	if err == nil && result != nil {
		if result.Status == "done" && len(result.Movies) > 0 {
			if !result.Contextual {
				moviesJSON, _ := json.Marshal(result.Movies)
				_ = s.repo.UpdateResult(ctx, taskID, string(moviesJSON))
			}
			return result, nil
		}
	}

	history, err := s.repo.GetHistoryByTaskID(ctx, taskID)
	if err == nil && history != nil && history.Result != "" {
		var movieIDs []string
		if err := json.Unmarshal([]byte(history.Result), &movieIDs); err == nil {
			return &redis.RecommendationResult{Status: "done", Movies: movieIDs}, nil
		}
	}
	return &redis.RecommendationResult{Status: "pending"}, nil
}

func (s *RecommendationService) GetHistory(ctx context.Context, userID string) ([]domain.RecommendationHistory, error) {
	return s.repo.GetHistory(ctx, userID)
}

func (s *RecommendationService) SaveRecommendation(ctx context.Context, userID, fromType, fromID, toType, toID string) (*domain.SavedRecommendation, error) {
	rec := &domain.SavedRecommendation{
		UserID:   userID,
		FromType: fromType,
		FromID:   fromID,
		ToType:   toType,
		ToID:     toID,
	}
	if err := s.repo.SaveRecommendation(ctx, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *RecommendationService) DeleteSavedRecommendation(ctx context.Context, id string) error {
	return s.repo.DeleteSavedRecommendation(ctx, id)
}

func (s *RecommendationService) GetSavedRecommendations(ctx context.Context, userID string) ([]domain.SavedRecommendation, error) {
	return s.repo.GetSavedRecommendations(ctx, userID)
}

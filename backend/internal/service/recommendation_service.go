package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/rabbitmq"
	"github.com/I000000/recly/internal/redis"
	"github.com/google/uuid"
)

type RecommendationService struct {
	repo           domain.RecommendationRepository
	publisher      rabbitmq.Publisher
	cache          redis.Cache
	sourceSelector *SourceSelector
}

func NewRecommendationService(
	repo domain.RecommendationRepository,
	pub rabbitmq.Publisher,
	cache redis.Cache,
	sourceSelector *SourceSelector,
) *RecommendationService {
	return &RecommendationService{
		repo:           repo,
		publisher:      pub,
		cache:          cache,
		sourceSelector: sourceSelector,
	}
}

func (s *RecommendationService) Request(
	ctx context.Context,
	userID string,
	selectedIDs []string,
	weights map[string]float64,
	excludeIDs []string,
	direction string,
	contextual bool,
) (string, error) {
	if len(selectedIDs) == 0 {
		var err error
		selectedIDs, weights, err = s.sourceSelector.Select(ctx, userID, nil)
		if err != nil {
			return "", err
		}
	} else if weights == nil {
		weights = make(map[string]float64, len(selectedIDs))
		for _, id := range selectedIDs {
			weights[id] = 1.0
		}
	}

	taskID := uuid.New().String()

	msg := rabbitmq.TaskMessage{
		TaskID:          taskID,
		UserID:          userID,
		SelectedIDs:     selectedIDs,
		SelectedWeights: weights,
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
		Contextual: contextual,
	}, 30*time.Minute); err != nil {
		// не фатально, но логируем
	}

	if !contextual {
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

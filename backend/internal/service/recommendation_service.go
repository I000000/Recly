package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/rabbitmq"
	"github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	log := logger.FromContext(ctx)
	log.Info("starting recommendation request",
		zap.String("user_id", userID),
		zap.String("direction", direction),
		zap.Int("selected_count", len(selectedIDs)),
		zap.Bool("contextual", contextual),
	)

	if len(selectedIDs) == 0 {
		var err error
		selectedIDs, weights, err = s.sourceSelector.Select(ctx, userID, nil)
		if err != nil {
			log.Error("failed to select sources", zap.Error(err), zap.String("user_id", userID))
			return "", err
		}
		log.Debug("selected sources from library", zap.Int("count", len(selectedIDs)))
	} else if weights == nil {
		weights = make(map[string]float64, len(selectedIDs))
		for _, id := range selectedIDs {
			weights[id] = 1.0
		}
		log.Debug("using explicit selected IDs with default weights", zap.Int("count", len(selectedIDs)))
	}

	taskID := uuid.New().String()
	log = log.With(zap.String("task_id", taskID))

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
		log.Error("failed to publish task to RabbitMQ", zap.Error(err))
		return "", err
	}
	log.Debug("task published to RabbitMQ")

	if err := s.cache.SetResult(ctx, taskID, redis.RecommendationResult{
		Status:     "pending",
		CreatedAt:  time.Now().Unix(),
		Contextual: contextual,
	}, 30*time.Minute); err != nil {
		log.Warn("failed to set result in Redis (non-critical)", zap.Error(err))
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
			log.Error("failed to save history", zap.Error(err))
			return "", err
		}
		log.Debug("history saved")
	}

	log.Info("recommendation request completed successfully", zap.String("task_id", taskID))
	return taskID, nil
}

func (s *RecommendationService) GetResult(ctx context.Context, taskID string) (*redis.RecommendationResult, error) {
	log := logger.FromContext(ctx).With(zap.String("task_id", taskID))
	log.Debug("fetching recommendation result")

	result, err := s.cache.GetResult(ctx, taskID)
	if err == nil && result != nil {
		if result.Status == "done" && len(result.Movies) > 0 {
			log.Debug("result found in cache", zap.Int("movies_count", len(result.Movies)))
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
			log.Debug("result found in history", zap.Int("movies_count", len(movieIDs)))
			return &redis.RecommendationResult{Status: "done", Movies: movieIDs}, nil
		}
	}

	log.Debug("result not found yet, returning pending")
	return &redis.RecommendationResult{Status: "pending"}, nil
}

func (s *RecommendationService) GetHistory(ctx context.Context, userID string) ([]domain.RecommendationHistory, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching recommendation history", zap.String("user_id", userID))
	history, err := s.repo.GetHistory(ctx, userID)
	if err != nil {
		log.Error("failed to get history", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("history fetched", zap.String("user_id", userID), zap.Int("count", len(history)))
	return history, nil
}

func (s *RecommendationService) SaveRecommendation(ctx context.Context, userID, fromType, fromID, toType, toID string) (*domain.SavedRecommendation, error) {
	log := logger.FromContext(ctx)
	log.Info("saving recommendation",
		zap.String("user_id", userID),
		zap.String("from_type", fromType),
		zap.String("from_id", fromID),
		zap.String("to_type", toType),
		zap.String("to_id", toID),
	)
	rec := &domain.SavedRecommendation{
		UserID:   userID,
		FromType: fromType,
		FromID:   fromID,
		ToType:   toType,
		ToID:     toID,
	}
	if err := s.repo.SaveRecommendation(ctx, rec); err != nil {
		log.Error("failed to save recommendation", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("recommendation saved", zap.String("user_id", userID), zap.String("saved_id", rec.ID))
	return rec, nil
}

func (s *RecommendationService) DeleteSavedRecommendation(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Info("deleting saved recommendation", zap.String("id", id))
	if err := s.repo.DeleteSavedRecommendation(ctx, id); err != nil {
		log.Error("failed to delete saved recommendation", zap.Error(err), zap.String("id", id))
		return err
	}
	log.Debug("saved recommendation deleted", zap.String("id", id))
	return nil
}

func (s *RecommendationService) GetSavedRecommendations(ctx context.Context, userID string) ([]domain.SavedRecommendation, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching saved recommendations", zap.String("user_id", userID))
	saved, err := s.repo.GetSavedRecommendations(ctx, userID)
	if err != nil {
		log.Error("failed to get saved recommendations", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("saved recommendations fetched", zap.String("user_id", userID), zap.Int("count", len(saved)))
	return saved, nil
}

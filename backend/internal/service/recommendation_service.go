package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

func (s *RecommendationService) Request(ctx context.Context, userID string, selectedIDs []string, direction string, weights map[string]float64) (string, error) {
	// Если список пуст – берём любимые книги или фильмы пользователя
	log.Println("DEBUG: Request called, libSvc is", s.libSvc)
	if len(selectedIDs) == 0 {
		switch direction {
		case "book_to_movie", "book_to_book":
			books, err := s.libSvc.GetBooks(ctx, userID)
			if err != nil {
				return "", err
			}
			for _, b := range books {
				selectedIDs = append(selectedIDs, b.BookID)
			}
		case "movie_to_book", "movie_to_movie":
			movies, err := s.libSvc.GetMovies(ctx, userID)
			if err != nil {
				return "", err
			}
			for _, m := range movies {
				selectedIDs = append(selectedIDs, m.MovieID)
			}
		}
		if len(selectedIDs) == 0 {
			return "", errors.New("no liked items to recommend from")
		}
	}

	taskID := uuid.New().String()
	// Публикуем задачу в RabbitMQ
	msg := rabbitmq.TaskMessage{
		TaskID:      taskID,
		UserID:      userID,
		SelectedIDs: selectedIDs,
		Direction:   direction,
		Weights:     weights,
	}
	if err := s.publisher.PublishRecommendationTask(ctx, msg); err != nil {
		return "", err
	}
	// Сохраняем начальный статус в кэше
	if err := s.cache.SetResult(ctx, taskID, redis.RecommendationResult{Status: "pending"}, 30*time.Minute); err != nil {
	}
	// Запись в историю БД
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
	// Пробуем Redis
	result, err := s.cache.GetResult(ctx, taskID)
	if err == nil && result != nil {
		return result, nil
	}
	// Если в Redis нет, ищем в истории БД
	history, err := s.repo.GetHistoryByTaskID(ctx, taskID)
	if err == nil && history != nil && history.Result != "" {
		// Парсим JSON-массив ID фильмов из поля result
		var movieIDs []string
		if err := json.Unmarshal([]byte(history.Result), &movieIDs); err == nil {
			return &redis.RecommendationResult{Status: "done", Movies: movieIDs}, nil
		}
	}
	// Если ничего не нашли, возвращаем pending
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

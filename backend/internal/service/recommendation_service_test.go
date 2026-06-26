package service

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecommendationService_Request_Success(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	pub := mocks.NewPublisher(t)
	cache := mocks.NewCache(t)

	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)
	repo.On("SaveHistory", mock.Anything, mock.AnythingOfType("*domain.RecommendationHistory")).Return(nil)

	svc := NewRecommendationService(repo, pub, cache, nil)

	taskID, err := svc.Request(
		context.Background(),
		"user-1",
		[]string{"book_1", "movie_1"},
		map[string]float64{"book_1": 1.0, "movie_1": 1.0},
		[]string{},
		"book_to_movie",
		false,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, taskID)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestRecommendationService_Request_PublishError(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	pub := mocks.NewPublisher(t)
	cache := mocks.NewCache(t)

	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(assert.AnError)

	svc := NewRecommendationService(repo, pub, cache, nil)

	taskID, err := svc.Request(
		context.Background(),
		"user-1",
		[]string{"book_1"},
		map[string]float64{"book_1": 1.0},
		[]string{},
		"book_to_movie",
		false,
	)

	assert.Error(t, err)
	assert.Empty(t, taskID)
	pub.AssertExpectations(t)
}

func TestRecommendationService_GetResult_FromCache(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	cache := mocks.NewCache(t)

	expected := &redis.RecommendationResult{
		Status:     "done",
		Movies:     []string{"movie_1", "movie_2"},
		Contextual: false,
	}
	cache.On("GetResult", mock.Anything, "task-1").Return(expected, nil)
	repo.On("UpdateResult", mock.Anything, "task-1", `["movie_1","movie_2"]`).Return(nil)

	svc := NewRecommendationService(repo, nil, cache, nil)

	result, err := svc.GetResult(context.Background(), "task-1")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestRecommendationService_GetResult_FromHistory(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	cache := mocks.NewCache(t)

	cache.On("GetResult", mock.Anything, "task-1").Return(nil, nil)
	history := &domain.RecommendationHistory{
		Result: `["movie_1", "movie_2"]`,
	}
	repo.On("GetHistoryByTaskID", mock.Anything, "task-1").Return(history, nil)

	svc := NewRecommendationService(repo, nil, cache, nil)

	result, err := svc.GetResult(context.Background(), "task-1")
	assert.NoError(t, err)
	assert.Equal(t, "done", result.Status)
	assert.Len(t, result.Movies, 2)
	repo.AssertExpectations(t)
}

func TestRecommendationService_Request_CacheSetError(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	pub := mocks.NewPublisher(t)
	cache := mocks.NewCache(t)

	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(assert.AnError)
	repo.On("SaveHistory", mock.Anything, mock.AnythingOfType("*domain.RecommendationHistory")).Return(nil)

	svc := NewRecommendationService(repo, pub, cache, nil)

	taskID, err := svc.Request(
		context.Background(),
		"user-1",
		[]string{"book_1"},
		map[string]float64{"book_1": 1.0},
		[]string{},
		"book_to_movie",
		false,
	)

	assert.NoError(t, err) // ошибка кэша не фатальна
	assert.NotEmpty(t, taskID)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestRecommendationService_Request_HistorySaveError(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	pub := mocks.NewPublisher(t)
	cache := mocks.NewCache(t)

	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)
	repo.On("SaveHistory", mock.Anything, mock.AnythingOfType("*domain.RecommendationHistory")).Return(assert.AnError)

	svc := NewRecommendationService(repo, pub, cache, nil)

	taskID, err := svc.Request(
		context.Background(),
		"user-1",
		[]string{"book_1"},
		map[string]float64{"book_1": 1.0},
		[]string{},
		"book_to_movie",
		false,
	)

	assert.Error(t, err) // ошибка сохранения истории должна прервать запрос
	assert.Empty(t, taskID)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestRecommendationService_Request_ContextualNoHistory(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	pub := mocks.NewPublisher(t)
	cache := mocks.NewCache(t)

	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	svc := NewRecommendationService(repo, pub, cache, nil)

	taskID, err := svc.Request(
		context.Background(),
		"user-1",
		[]string{"book_1"},
		map[string]float64{"book_1": 1.0},
		[]string{},
		"book_to_movie",
		true, // contextual
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, taskID)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestRecommendationService_GetResult_UpdateResultError(t *testing.T) {
	repo := mocks.NewRecommendationRepository(t)
	cache := mocks.NewCache(t)

	expected := &redis.RecommendationResult{
		Status:     "done",
		Movies:     []string{"movie_1"},
		Contextual: false,
	}
	cache.On("GetResult", mock.Anything, "task-1").Return(expected, nil)
	repo.On("UpdateResult", mock.Anything, "task-1", `["movie_1"]`).Return(assert.AnError)

	svc := NewRecommendationService(repo, nil, cache, nil)

	result, err := svc.GetResult(context.Background(), "task-1")
	assert.NoError(t, err) // ошибка UpdateResult не фатальна
	assert.Equal(t, expected, result)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

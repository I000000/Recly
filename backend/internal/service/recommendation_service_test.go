package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/I000000/recly/internal/rabbitmq"
	redisPkg "github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/internal/service"
	"github.com/I000000/recly/mocks"
)

// Ручные моки для Publisher и Cache
type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) PublishRecommendationTask(ctx context.Context, msg rabbitmq.TaskMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockPublisher) Close() error {
	return nil
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) SetResult(ctx context.Context, taskID string, result redisPkg.RecommendationResult, ttl time.Duration) error {
	args := m.Called(ctx, taskID, result, ttl)
	return args.Error(0)
}

func (m *mockCache) GetResult(ctx context.Context, taskID string) (*redisPkg.RecommendationResult, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*redisPkg.RecommendationResult), args.Error(1)
}

func TestRecommendationService_Request_PublishesAndCaches(t *testing.T) {
	recRepo := &mocks.RecommendationRepository{}
	pub := new(mockPublisher)
	cache := new(mockCache)

	recRepo.On("SaveHistory", mock.Anything, mock.AnythingOfType("*domain.RecommendationHistory")).Return(nil)
	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("redis.RecommendationResult"), mock.Anything).Return(nil)

	svc := service.NewRecommendationService(recRepo, pub, cache)
	taskID, err := svc.Request(context.Background(), "user1", []string{"book1"}, "book_to_movie", map[string]float64{"text": 0.6})

	assert.NoError(t, err)
	assert.NotEmpty(t, taskID)
	recRepo.AssertExpectations(t)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
}

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/I000000/recly/internal/handler"
	"github.com/I000000/recly/internal/rabbitmq"
	redisPkg "github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/internal/service"
	"github.com/I000000/recly/mocks"
)

// ─── Ручные моки для Publisher и Cache ───

type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) PublishRecommendationTask(ctx context.Context, msg rabbitmq.TaskMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockPublisher) Close() error { return nil }

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

// ─── Тесты ───

func TestRecommendationHandler_Request(t *testing.T) {
	recRepo := &mocks.RecommendationRepository{}
	pub := new(mockPublisher)
	cache := new(mockCache)

	// Настраиваем моки
	recRepo.On("SaveHistory", mock.Anything, mock.AnythingOfType("*domain.RecommendationHistory")).Return(nil)
	pub.On("PublishRecommendationTask", mock.Anything, mock.AnythingOfType("rabbitmq.TaskMessage")).Return(nil)
	cache.On("SetResult", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("redis.RecommendationResult"), mock.Anything).Return(nil)

	// Создаём сервис с тремя зависимостями
	svc := service.NewRecommendationService(recRepo, pub, cache)
	h := handler.NewRecommendationHandler(svc)

	r := gin.New()
	r.POST("/recommend", func(c *gin.Context) { c.Set("user_id", "user1"); h.Request(c) })

	body := `{"selected_ids":["book1"],"direction":"book_to_movie"}`
	req, _ := http.NewRequest("POST", "/recommend", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "pending", resp["status"])
	assert.NotEmpty(t, resp["task_id"])
	recRepo.AssertExpectations(t)
	pub.AssertExpectations(t)
	cache.AssertExpectations(t)
}

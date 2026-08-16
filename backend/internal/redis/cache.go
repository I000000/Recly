//go:generate mockery --name Cache --output ../../mocks --outpkg mocks --case underscore
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/I000000/recly/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RecommendationResult struct {
	Status     string   `json:"status"`
	Movies     []string `json:"movies,omitempty"`
	Error      string   `json:"error,omitempty"`
	CreatedAt  int64    `json:"created_at,omitempty"`
	Contextual bool     `json:"contextual,omitempty"`
}

type Cache interface {
	SetResult(ctx context.Context, taskID string, result RecommendationResult, ttl time.Duration) error
	GetResult(ctx context.Context, taskID string) (*RecommendationResult, error)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: rdb}
}

func (r *RedisCache) SetResult(ctx context.Context, taskID string, result RecommendationResult, ttl time.Duration) error {
	log := logger.FromContext(ctx)
	data, err := json.Marshal(result)
	if err != nil {
		log.Error("failed to marshal result for Redis", zap.Error(err), zap.String("task_id", taskID))
		return fmt.Errorf("marshal: %w", err)
	}
	err = r.client.Set(ctx, "rec:"+taskID, data, ttl).Err()
	if err != nil {
		log.Error("failed to set result in Redis", zap.Error(err), zap.String("task_id", taskID))
		return err
	}
	log.Debug("result stored in Redis", zap.String("task_id", taskID))
	return nil
}

func (r *RedisCache) GetResult(ctx context.Context, taskID string) (*RecommendationResult, error) {
	log := logger.FromContext(ctx)
	val, err := r.client.Get(ctx, "rec:"+taskID).Bytes()
	if err == redis.Nil {
		log.Debug("result not found in Redis", zap.String("task_id", taskID))
		return nil, nil
	}
	if err != nil {
		log.Error("failed to get result from Redis", zap.Error(err), zap.String("task_id", taskID))
		return nil, fmt.Errorf("redis get: %w", err)
	}
	var res RecommendationResult
	if err := json.Unmarshal(val, &res); err != nil {
		log.Error("failed to unmarshal result from Redis", zap.Error(err), zap.String("task_id", taskID))
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	log.Debug("result retrieved from Redis", zap.String("task_id", taskID))
	return &res, nil
}

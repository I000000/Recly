//go:generate mockery --name Cache --output ../../mocks --outpkg mocks --case underscore
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return r.client.Set(ctx, "rec:"+taskID, data, ttl).Err()
}

func (r *RedisCache) GetResult(ctx context.Context, taskID string) (*RecommendationResult, error) {
	val, err := r.client.Get(ctx, "rec:"+taskID).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}
	var res RecommendationResult
	if err := json.Unmarshal(val, &res); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &res, nil
}

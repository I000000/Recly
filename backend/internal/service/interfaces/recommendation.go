//go:generate mockery --name RecommendationService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/redis"
)

type RecommendationService interface {
	Request(ctx context.Context, userID string, selectedIDs []string, weights map[string]float64, excludeIDs []string, direction string, contextual bool) (string, error)
	GetResult(ctx context.Context, taskID string) (*redis.RecommendationResult, error)
	GetHistory(ctx context.Context, userID string) ([]domain.RecommendationHistory, error)
	SaveRecommendation(ctx context.Context, userID, fromType, fromID, toType, toID string) (*domain.SavedRecommendation, error)
	DeleteSavedRecommendation(ctx context.Context, id string) error
	GetSavedRecommendations(ctx context.Context, userID string) ([]domain.SavedRecommendation, error)
}

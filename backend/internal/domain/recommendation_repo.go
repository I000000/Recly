package domain

import "context"

type RecommendationRepository interface {
	SaveHistory(ctx context.Context, entry *RecommendationHistory) error
	GetHistory(ctx context.Context, userID string) ([]RecommendationHistory, error)
	GetHistoryByTaskID(ctx context.Context, taskID string) (*RecommendationHistory, error)
	UpdateResult(ctx context.Context, taskID string, resultJSON string) error
	SaveRecommendation(ctx context.Context, rec *SavedRecommendation) error
	DeleteSavedRecommendation(ctx context.Context, id string) error
	GetSavedRecommendations(ctx context.Context, userID string) ([]SavedRecommendation, error)
}

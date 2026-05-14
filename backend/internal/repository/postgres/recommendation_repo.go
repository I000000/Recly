package postgres

import (
	"context"
	"encoding/json"

	"github.com/I000000/recly/internal/domain"
)

type RecommendationRepo struct {
	pool Querier
}

func NewRecommendationRepo(pool Querier) domain.RecommendationRepository {
	return &RecommendationRepo{pool: pool}
}

func (r *RecommendationRepo) SaveHistory(ctx context.Context, entry *domain.RecommendationHistory) error {
	selJSON, err := json.Marshal(entry.SelectedIDs)
	if err != nil {
		return err
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_recommendation_history (user_id, task_id, selected_ids, direction, weights, result)
         VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		entry.UserID, entry.TaskID, string(selJSON), entry.Direction, entry.Weights, entry.Result,
	).Scan(&entry.ID, &entry.CreatedAt)
}

func (r *RecommendationRepo) GetHistory(ctx context.Context, userID string) ([]domain.RecommendationHistory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, task_id, selected_ids, direction, weights, result, created_at
         FROM user_recommendation_history WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.RecommendationHistory
	for rows.Next() {
		var e domain.RecommendationHistory
		var selStr string
		if err := rows.Scan(&e.ID, &e.UserID, &e.TaskID, &selStr, &e.Direction, &e.Weights, &e.Result, &e.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(selStr), &e.SelectedIDs)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *RecommendationRepo) SaveRecommendation(ctx context.Context, rec *domain.SavedRecommendation) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO user_saved_recommendations (user_id, from_type, from_id, to_type, to_id)
         VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, from_type, from_id, to_type, to_id) DO NOTHING
         RETURNING id, saved_at`,
		rec.UserID, rec.FromType, rec.FromID, rec.ToType, rec.ToID,
	).Scan(&rec.ID, &rec.SavedAt)
}

func (r *RecommendationRepo) DeleteSavedRecommendation(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM user_saved_recommendations WHERE id = $1`, id)
	return err
}

func (r *RecommendationRepo) GetSavedRecommendations(ctx context.Context, userID string) ([]domain.SavedRecommendation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, from_type, from_id, to_type, to_id, saved_at
         FROM user_saved_recommendations WHERE user_id = $1 ORDER BY saved_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recs []domain.SavedRecommendation
	for rows.Next() {
		var r domain.SavedRecommendation
		if err := rows.Scan(&r.ID, &r.UserID, &r.FromType, &r.FromID, &r.ToType, &r.ToID, &r.SavedAt); err != nil {
			return nil, err
		}
		recs = append(recs, r)
	}
	return recs, rows.Err()
}

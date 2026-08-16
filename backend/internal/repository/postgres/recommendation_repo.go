package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type RecommendationRepo struct {
	pool Querier
}

func NewRecommendationRepo(pool Querier) domain.RecommendationRepository {
	return &RecommendationRepo{pool: pool}
}

func (r *RecommendationRepo) SaveHistory(ctx context.Context, entry *domain.RecommendationHistory) error {
	log := logger.FromContext(ctx)
	log.Debug("saving recommendation history", zap.String("task_id", entry.TaskID), zap.String("user_id", entry.UserID))
	selJSON, err := json.Marshal(entry.SelectedIDs)
	if err != nil {
		log.Error("failed to marshal selected IDs", zap.Error(err))
		return err
	}
	var result interface{}
	if entry.Result == "" {
		result = nil
	} else {
		result = entry.Result
	}
	err = r.pool.QueryRow(ctx,
		`INSERT INTO user_recommendation_history (user_id, task_id, selected_ids, direction, weights, result)
         VALUES ($1, $2, $3, $4, $5::jsonb, $6) RETURNING id, created_at`,
		entry.UserID, entry.TaskID, string(selJSON), entry.Direction, entry.Weights, result,
	).Scan(&entry.ID, &entry.CreatedAt)
	if err != nil {
		log.Error("failed to save recommendation history", zap.Error(err), zap.String("task_id", entry.TaskID))
	}
	return err
}

func (r *RecommendationRepo) UpdateResult(ctx context.Context, taskID string, resultJSON string) error {
	log := logger.FromContext(ctx)
	log.Debug("updating recommendation result", zap.String("task_id", taskID))
	_, err := r.pool.Exec(ctx,
		`UPDATE user_recommendation_history SET result = $1::jsonb WHERE task_id = $2`,
		resultJSON, taskID)
	if err != nil {
		log.Error("failed to update recommendation result", zap.Error(err), zap.String("task_id", taskID))
	}
	return err
}

func (r *RecommendationRepo) GetHistory(ctx context.Context, userID string) ([]domain.RecommendationHistory, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting recommendation history", zap.String("user_id", userID))
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, task_id, selected_ids, direction, weights, result, created_at
         FROM user_recommendation_history WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Error("failed to query recommendation history", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var entries []domain.RecommendationHistory
	for rows.Next() {
		var e domain.RecommendationHistory
		var selStr string
		var resultStr sql.NullString
		if err := rows.Scan(&e.ID, &e.UserID, &e.TaskID, &selStr, &e.Direction, &e.Weights, &resultStr, &e.CreatedAt); err != nil {
			log.Error("failed to scan recommendation history", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		e.Result = resultStr.String
		json.Unmarshal([]byte(selStr), &e.SelectedIDs)
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("recommendation history retrieved", zap.String("user_id", userID), zap.Int("count", len(entries)))
	return entries, nil
}

func (r *RecommendationRepo) GetHistoryByTaskID(ctx context.Context, taskID string) (*domain.RecommendationHistory, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting recommendation history by task ID", zap.String("task_id", taskID))
	var e domain.RecommendationHistory
	var selStr string
	var resultStr sql.NullString
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, task_id, selected_ids, direction, weights, result, created_at
         FROM user_recommendation_history WHERE task_id = $1`,
		taskID,
	).Scan(&e.ID, &e.UserID, &e.TaskID, &selStr, &e.Direction, &e.Weights, &resultStr, &e.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("history not found by task ID", zap.String("task_id", taskID))
			return nil, nil
		}
		log.Error("failed to get history by task ID", zap.Error(err), zap.String("task_id", taskID))
		return nil, err
	}
	e.Result = resultStr.String
	json.Unmarshal([]byte(selStr), &e.SelectedIDs)
	log.Debug("history found by task ID", zap.String("task_id", taskID))
	return &e, nil
}

func (r *RecommendationRepo) SaveRecommendation(ctx context.Context, rec *domain.SavedRecommendation) error {
	log := logger.FromContext(ctx)
	log.Debug("saving recommendation", zap.String("user_id", rec.UserID), zap.String("from_id", rec.FromID), zap.String("to_id", rec.ToID))
	err := r.pool.QueryRow(ctx,
		`INSERT INTO user_saved_recommendations (user_id, from_type, from_id, to_type, to_id)
         VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, from_type, from_id, to_type, to_id) DO NOTHING
         RETURNING id, saved_at`,
		rec.UserID, rec.FromType, rec.FromID, rec.ToType, rec.ToID,
	).Scan(&rec.ID, &rec.SavedAt)
	if err != nil {
		log.Error("failed to save recommendation", zap.Error(err), zap.String("user_id", rec.UserID))
	}
	return err
}

func (r *RecommendationRepo) DeleteSavedRecommendation(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Debug("deleting saved recommendation", zap.String("id", id))
	_, err := r.pool.Exec(ctx, `DELETE FROM user_saved_recommendations WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete saved recommendation", zap.Error(err), zap.String("id", id))
	}
	return err
}

func (r *RecommendationRepo) GetSavedRecommendations(ctx context.Context, userID string) ([]domain.SavedRecommendation, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting saved recommendations", zap.String("user_id", userID))
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, from_type, from_id, to_type, to_id, saved_at
         FROM user_saved_recommendations WHERE user_id = $1 ORDER BY saved_at DESC`, userID)
	if err != nil {
		log.Error("failed to query saved recommendations", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var recs []domain.SavedRecommendation
	for rows.Next() {
		var r domain.SavedRecommendation
		if err := rows.Scan(&r.ID, &r.UserID, &r.FromType, &r.FromID, &r.ToType, &r.ToID, &r.SavedAt); err != nil {
			log.Error("failed to scan saved recommendation", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		recs = append(recs, r)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("saved recommendations retrieved", zap.String("user_id", userID), zap.Int("count", len(recs)))
	return recs, nil
}

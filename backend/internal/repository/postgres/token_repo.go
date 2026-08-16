package postgres

import (
	"context"
	"errors"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type TokenRepo struct {
	pool Querier
}

func NewTokenRepo(pool Querier) domain.TokenRepository {
	return &TokenRepo{pool: pool}
}

func (r *TokenRepo) StoreRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	log := logger.FromContext(ctx)
	log.Debug("storing refresh token", zap.String("user_id", rt.UserID))
	err := r.pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at`,
		rt.UserID, rt.TokenHash, rt.ExpiresAt,
	).Scan(&rt.ID, &rt.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn("duplicate token", zap.String("user_id", rt.UserID))
			return domain.ErrDuplicateEmail
		}
		log.Error("failed to store refresh token", zap.Error(err), zap.String("user_id", rt.UserID))
		return err
	}
	log.Debug("refresh token stored", zap.String("token_id", rt.ID))
	return nil
}

func (r *TokenRepo) GetRefreshToken(ctx context.Context, id string) (*domain.RefreshToken, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting refresh token", zap.String("token_id", id))
	rt := &domain.RefreshToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE id = $1`,
		id,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("refresh token not found", zap.String("token_id", id))
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get refresh token", zap.Error(err), zap.String("token_id", id))
		return nil, err
	}
	log.Debug("refresh token found", zap.String("token_id", rt.ID))
	return rt, nil
}

func (r *TokenRepo) DeleteRefreshToken(ctx context.Context, id string) error {
	log := logger.FromContext(ctx)
	log.Debug("deleting refresh token", zap.String("token_id", id))
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	if err != nil {
		log.Error("failed to delete refresh token", zap.Error(err), zap.String("token_id", id))
	}
	return err
}

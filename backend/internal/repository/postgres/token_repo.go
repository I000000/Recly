package postgres

import (
	"context"
	"errors"

	"github.com/I000000/recly/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TokenRepo struct {
	pool Querier
}

func NewTokenRepo(pool Querier) domain.TokenRepository {
	return &TokenRepo{pool: pool}
}

func (r *TokenRepo) StoreRefreshToken(ctx context.Context, rt *domain.RefreshToken) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at`,
		rt.UserID, rt.TokenHash, rt.ExpiresAt,
	).Scan(&rt.ID, &rt.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *TokenRepo) GetRefreshToken(ctx context.Context, id string) (*domain.RefreshToken, error) {
	rt := &domain.RefreshToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE id = $1`,
		id,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return rt, nil
}

func (r *TokenRepo) DeleteRefreshToken(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	return err
}

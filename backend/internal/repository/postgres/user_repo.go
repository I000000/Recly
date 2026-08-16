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

type UserRepo struct {
	pool Querier
}

func NewUserRepo(pool Querier) domain.UserRepository {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	log := logger.FromContext(ctx)
	log.Debug("creating user", zap.String("email", user.Email))
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at`,
		user.Email, user.PasswordHash, user.Name,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn("duplicate email", zap.String("email", user.Email))
			return domain.ErrDuplicateEmail
		}
		log.Error("failed to create user", zap.Error(err), zap.String("email", user.Email))
		return err
	}
	log.Debug("user created", zap.String("user_id", user.ID))
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting user by email", zap.String("email", email))
	user := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, created_at, onboarding_completed, avatar_url FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.OnboardingCompleted, &user.AvatarURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("user not found by email", zap.String("email", email))
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	log.Debug("user found", zap.String("user_id", user.ID))
	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting user by ID", zap.String("user_id", id))
	user := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, name, created_at, onboarding_completed, avatar_url FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.OnboardingCompleted, &user.AvatarURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("user not found by ID", zap.String("user_id", id))
			return nil, domain.ErrNotFound
		}
		log.Error("failed to get user by ID", zap.Error(err), zap.String("user_id", id))
		return nil, err
	}
	log.Debug("user found", zap.String("user_id", user.ID))
	return user, nil
}

func (r *UserRepo) UpdateOnboardingCompleted(ctx context.Context, userID string, completed bool) error {
	log := logger.FromContext(ctx)
	log.Debug("updating onboarding status", zap.String("user_id", userID), zap.Bool("completed", completed))
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET onboarding_completed = $1 WHERE id = $2`,
		completed, userID,
	)
	if err != nil {
		log.Error("failed to update onboarding status", zap.Error(err), zap.String("user_id", userID))
		return err
	}
	log.Debug("onboarding status updated", zap.String("user_id", userID))
	return nil
}

func (r *UserRepo) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	log := logger.FromContext(ctx)
	log.Debug("updating avatar URL", zap.String("user_id", userID), zap.String("avatar_url", avatarURL))
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET avatar_url = $1 WHERE id = $2`,
		avatarURL, userID,
	)
	if err != nil {
		log.Error("failed to update avatar", zap.Error(err), zap.String("user_id", userID))
		return err
	}
	log.Debug("avatar URL updated", zap.String("user_id", userID))
	return nil
}

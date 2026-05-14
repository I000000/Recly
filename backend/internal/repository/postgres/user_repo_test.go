package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/repository/postgres"
)

func TestUserRepo_Create_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	repo := postgres.NewUserRepo(mockPool)
	user := &domain.User{
		Email:        "test@example.com",
		PasswordHash: "hashed",
		Name:         "Test",
	}

	mockPool.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Email, user.PasswordHash, user.Name).
		WillReturnRows(
			mockPool.NewRows([]string{"id", "created_at"}).
				AddRow("550e8400-e29b-41d4-a716-446655440000", time.Now()),
		)

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", user.ID)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

func TestUserRepo_Create_DuplicateEmail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	repo := postgres.NewUserRepo(mockPool)
	user := &domain.User{
		Email:        "dup@example.com",
		PasswordHash: "hashed",
		Name:         "Dup",
	}

	mockPool.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Email, user.PasswordHash, user.Name).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err = repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateEmail, err)
	assert.NoError(t, mockPool.ExpectationsWereMet())
}

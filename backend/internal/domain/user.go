package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

type User struct {
	ID                  string    `json:"id"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"-"`
	Name                string    `json:"name"`
	AvatarURL           string    `db:"avatar_url" json:"avatar_url"`
	OnboardingCompleted bool      `db:"onboarding_completed" json:"onboarding_completed"`
	CreatedAt           time.Time `json:"created_at"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateOnboardingCompleted(ctx context.Context, userID string, completed bool) error
	UpdateAvatar(ctx context.Context, userID, avatarURL string) error
}

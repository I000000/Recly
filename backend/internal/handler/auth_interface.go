package handler

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

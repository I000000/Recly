//go:generate mockery --name AuthService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

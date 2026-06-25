package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/I000000/recly/internal/domain"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	CompleteOnboarding(ctx context.Context, userID string) error
	UpdateAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error)
}

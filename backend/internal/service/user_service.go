package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) CompleteOnboarding(ctx context.Context, userID string) error {
	return s.repo.UpdateOnboardingCompleted(ctx, userID, true)
}

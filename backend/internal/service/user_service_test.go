package service

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetUserByID_Success(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	expected := &domain.User{ID: "user-1", Email: "test@example.com"}
	repo.On("GetByID", mock.Anything, "user-1").Return(expected, nil)

	svc := NewUserService(repo)
	user, err := svc.GetUserByID(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, expected, user)
	repo.AssertExpectations(t)
}

func TestUserService_CompleteOnboarding_Success(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	repo.On("UpdateOnboardingCompleted", mock.Anything, "user-1", true).Return(nil)

	svc := NewUserService(repo)
	err := svc.CompleteOnboarding(context.Background(), "user-1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	repo.On("GetByID", mock.Anything, "user-1").Return(nil, assert.AnError)

	svc := NewUserService(repo)
	user, err := svc.GetUserByID(context.Background(), "user-1")
	assert.Error(t, err)
	assert.Nil(t, user)
	repo.AssertExpectations(t)
}

func TestUserService_CompleteOnboarding_Error(t *testing.T) {
	repo := mocks.NewUserRepository(t)
	repo.On("UpdateOnboardingCompleted", mock.Anything, "user-1", true).Return(assert.AnError)

	svc := NewUserService(repo)
	err := svc.CompleteOnboarding(context.Background(), "user-1")
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

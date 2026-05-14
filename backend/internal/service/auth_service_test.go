package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/I000000/recly/internal/service"
	"github.com/I000000/recly/mocks"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := &mocks.UserRepository{}
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	mockTokenRepo := new(mocks.TokenRepository)

	svc := service.NewAuthService(mockUserRepo, mockTokenRepo, "secret", 15, 43200)
	user, err := svc.Register(context.Background(), "test@example.com", "12345678", "Test")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	mockUserRepo.AssertExpectations(t)
}

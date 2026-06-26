package service

import (
	"context"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	tokenRepo := mocks.NewTokenRepository(t)

	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	svc := NewAuthService(userRepo, tokenRepo, "test-secret", 60, 43200)

	user, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
	assert.NoError(t, err)

	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	tokenRepo := mocks.NewTokenRepository(t)

	userRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(domain.ErrDuplicateEmail)

	svc := NewAuthService(userRepo, tokenRepo, "test-secret", 60, 43200)

	user, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrDuplicateEmail, err)

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	tokenRepo := mocks.NewTokenRepository(t)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: string(hashed),
	}

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
	tokenRepo.On("StoreRefreshToken", mock.Anything, mock.AnythingOfType("*domain.RefreshToken")).Return(nil)

	svc := NewAuthService(userRepo, tokenRepo, "test-secret", 60, 43200)

	access, refresh, err := svc.Login(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	tokenRepo := mocks.NewTokenRepository(t)

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domain.ErrNotFound)

	svc := NewAuthService(userRepo, tokenRepo, "test-secret", 60, 43200)

	access, refresh, err := svc.Login(context.Background(), "test@example.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, "invalid credentials", err.Error())

	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	userRepo := mocks.NewUserRepository(t)
	tokenRepo := mocks.NewTokenRepository(t)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           "user-1",
		Email:        "test@example.com",
		PasswordHash: string(hashed),
	}

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

	svc := NewAuthService(userRepo, tokenRepo, "test-secret", 60, 43200)

	access, refresh, err := svc.Login(context.Background(), "test@example.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, "invalid credentials", err.Error())

	userRepo.AssertExpectations(t)
}

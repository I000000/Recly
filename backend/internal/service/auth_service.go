package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwtSecret  string
	accessTTL  int
	refreshTTL int
}

func NewAuthService(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	jwtSecret string,
	accessTTL, refreshTTL int,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		Name:         name,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	accessToken, err := jwt.GenerateAccessToken(user.ID, s.jwtSecret, s.accessTTL)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenPlain := hex.EncodeToString(b)
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenPlain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: string(tokenHash),
		ExpiresAt: time.Now().Add(time.Duration(s.refreshTTL) * time.Minute),
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, rt); err != nil {
		return "", err
	}
	return rt.ID + ":" + tokenPlain, nil
}

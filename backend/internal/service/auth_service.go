package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/jwt"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
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
	log := logger.FromContext(ctx)
	log.Info("registering user", zap.String("email", email))

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))
		return nil, err
	}
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashed),
		Name:         name,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			log.Warn("registration failed: duplicate email", zap.String("email", email))
		} else {
			log.Error("registration failed: db error", zap.Error(err), zap.String("email", email))
		}
		return nil, err
	}
	log.Info("user registered successfully", zap.String("user_id", user.ID), zap.String("email", email))
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	log := logger.FromContext(ctx)
	log.Info("login attempt", zap.String("email", email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Warn("login failed: user not found", zap.String("email", email))
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn("login failed: wrong password", zap.String("email", email))
		return "", "", errors.New("invalid credentials")
	}
	accessToken, err := jwt.GenerateAccessToken(user.ID, s.jwtSecret, s.accessTTL)
	if err != nil {
		log.Error("failed to generate access token", zap.Error(err), zap.String("user_id", user.ID))
		return "", "", err
	}
	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		log.Error("failed to generate refresh token", zap.Error(err), zap.String("user_id", user.ID))
		return "", "", err
	}
	log.Info("login successful", zap.String("user_id", user.ID), zap.String("email", email))
	return accessToken, refreshToken, nil
}

func (s *AuthService) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	log := logger.FromContext(ctx)
	log.Debug("generating refresh token", zap.String("user_id", userID))

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Error("failed to read random bytes", zap.Error(err))
		return "", err
	}
	tokenPlain := hex.EncodeToString(b)
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenPlain), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash refresh token", zap.Error(err))
		return "", err
	}
	rt := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: string(tokenHash),
		ExpiresAt: time.Now().Add(time.Duration(s.refreshTTL) * time.Minute),
	}
	if err := s.tokenRepo.StoreRefreshToken(ctx, rt); err != nil {
		log.Error("failed to store refresh token", zap.Error(err), zap.String("user_id", userID))
		return "", err
	}
	log.Debug("refresh token stored successfully", zap.String("token_id", rt.ID))
	return rt.ID + ":" + tokenPlain, nil
}

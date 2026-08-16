package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching user by ID", zap.String("user_id", id))
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Warn("user not found", zap.Error(err), zap.String("user_id", id))
		return nil, err
	}
	return user, nil
}

func (s *UserService) CompleteOnboarding(ctx context.Context, userID string) error {
	log := logger.FromContext(ctx)
	log.Info("completing onboarding", zap.String("user_id", userID))
	if err := s.repo.UpdateOnboardingCompleted(ctx, userID, true); err != nil {
		log.Error("failed to complete onboarding", zap.Error(err), zap.String("user_id", userID))
		return err
	}
	log.Debug("onboarding completed", zap.String("user_id", userID))
	return nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	log := logger.FromContext(ctx)
	log.Info("updating avatar", zap.String("user_id", userID), zap.String("filename", header.Filename), zap.Int64("size", header.Size))

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%d%s", userID, time.Now().UnixNano(), ext)
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Error("failed to create upload dir", zap.Error(err), zap.String("dir", uploadDir))
		return "", err
	}
	path := filepath.Join(uploadDir, filename)
	dst, err := os.Create(path)
	if err != nil {
		log.Error("failed to create file", zap.Error(err), zap.String("path", path))
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(path)
		log.Error("failed to copy file", zap.Error(err), zap.String("path", path))
		return "", err
	}
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		os.Remove(path)
		log.Error("failed to update avatar in DB", zap.Error(err), zap.String("user_id", userID))
		return "", err
	}
	log.Info("avatar updated", zap.String("user_id", userID), zap.String("avatar_url", avatarURL))
	return avatarURL, nil
}

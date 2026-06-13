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

func (s *UserService) UpdateAvatar(ctx context.Context, userID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%s_%d%s", userID, time.Now().UnixNano(), ext)
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(uploadDir, filename)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(path)
		return "", err
	}
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := s.repo.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		os.Remove(path)
		return "", err
	}
	return avatarURL, nil
}

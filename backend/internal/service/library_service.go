package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type LibraryService struct {
	repo domain.LibraryRepository
}

func NewLibraryService(repo domain.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (s *LibraryService) AddBook(ctx context.Context, userID, bookID string) error {
	return s.repo.AddLikedBook(ctx, userID, bookID)
}

func (s *LibraryService) RemoveBook(ctx context.Context, userID, bookID string) error {
	return s.repo.RemoveLikedBook(ctx, userID, bookID)
}

func (s *LibraryService) GetBooks(ctx context.Context, userID string) ([]domain.LikedBook, error) {
	return s.repo.GetLikedBooks(ctx, userID)
}

func (s *LibraryService) AddMovie(ctx context.Context, userID, movieID string) error {
	return s.repo.AddLikedMovie(ctx, userID, movieID)
}

func (s *LibraryService) RemoveMovie(ctx context.Context, userID, movieID string) error {
	return s.repo.RemoveLikedMovie(ctx, userID, movieID)
}

func (s *LibraryService) GetMovies(ctx context.Context, userID string) ([]domain.LikedMovie, error) {
	return s.repo.GetLikedMovies(ctx, userID)
}

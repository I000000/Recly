package service

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type LibraryService struct {
	repo domain.LibraryRepository
}

func NewLibraryService(repo domain.LibraryRepository) *LibraryService {
	return &LibraryService{repo: repo}
}

func (s *LibraryService) AddBook(ctx context.Context, userID, bookID string) error {
	log := logger.FromContext(ctx)
	log.Info("adding book to library", zap.String("user_id", userID), zap.String("book_id", bookID))
	if err := s.repo.AddLikedBook(ctx, userID, bookID); err != nil {
		log.Error("failed to add book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
		return err
	}
	log.Debug("book added", zap.String("user_id", userID), zap.String("book_id", bookID))
	return nil
}

func (s *LibraryService) RemoveBook(ctx context.Context, userID, bookID string) error {
	log := logger.FromContext(ctx)
	log.Info("removing book from library", zap.String("user_id", userID), zap.String("book_id", bookID))
	if err := s.repo.RemoveLikedBook(ctx, userID, bookID); err != nil {
		log.Error("failed to remove book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
		return err
	}
	log.Debug("book removed", zap.String("user_id", userID), zap.String("book_id", bookID))
	return nil
}

func (s *LibraryService) GetBooks(ctx context.Context, userID string) ([]domain.LikedBook, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching liked books", zap.String("user_id", userID))
	books, err := s.repo.GetLikedBooks(ctx, userID)
	if err != nil {
		log.Error("failed to get liked books", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("liked books fetched", zap.String("user_id", userID), zap.Int("count", len(books)))
	return books, nil
}

func (s *LibraryService) AddMovie(ctx context.Context, userID, movieID string) error {
	log := logger.FromContext(ctx)
	log.Info("adding movie to library", zap.String("user_id", userID), zap.String("movie_id", movieID))
	if err := s.repo.AddLikedMovie(ctx, userID, movieID); err != nil {
		log.Error("failed to add movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
		return err
	}
	log.Debug("movie added", zap.String("user_id", userID), zap.String("movie_id", movieID))
	return nil
}

func (s *LibraryService) RemoveMovie(ctx context.Context, userID, movieID string) error {
	log := logger.FromContext(ctx)
	log.Info("removing movie from library", zap.String("user_id", userID), zap.String("movie_id", movieID))
	if err := s.repo.RemoveLikedMovie(ctx, userID, movieID); err != nil {
		log.Error("failed to remove movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
		return err
	}
	log.Debug("movie removed", zap.String("user_id", userID), zap.String("movie_id", movieID))
	return nil
}

func (s *LibraryService) GetMovies(ctx context.Context, userID string) ([]domain.LikedMovie, error) {
	log := logger.FromContext(ctx)
	log.Debug("fetching liked movies", zap.String("user_id", userID))
	movies, err := s.repo.GetLikedMovies(ctx, userID)
	if err != nil {
		log.Error("failed to get liked movies", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("liked movies fetched", zap.String("user_id", userID), zap.Int("count", len(movies)))
	return movies, nil
}

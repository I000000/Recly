package domain

import "context"

type LibraryRepository interface {
	AddLikedBook(ctx context.Context, userID, bookID string) error
	RemoveLikedBook(ctx context.Context, userID, bookID string) error
	GetLikedBooks(ctx context.Context, userID string) ([]LikedBook, error)
	AddLikedMovie(ctx context.Context, userID, movieID string) error
	RemoveLikedMovie(ctx context.Context, userID, movieID string) error
	GetLikedMovies(ctx context.Context, userID string) ([]LikedMovie, error)
}

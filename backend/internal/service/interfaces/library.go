//go:generate mockery --name LibraryService --output ../../../mocks --outpkg mocks --case underscore
package interfaces

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type LibraryService interface {
	AddBook(ctx context.Context, userID, bookID string) error
	RemoveBook(ctx context.Context, userID, bookID string) error
	GetBooks(ctx context.Context, userID string) ([]domain.LikedBook, error)
	AddMovie(ctx context.Context, userID, movieID string) error
	RemoveMovie(ctx context.Context, userID, movieID string) error
	GetMovies(ctx context.Context, userID string) ([]domain.LikedMovie, error)
}

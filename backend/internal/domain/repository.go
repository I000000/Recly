package domain

import "context"

//go:generate mockery --name UserRepository --output ../../mocks --outpkg mocks --case underscore
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateOnboardingCompleted(ctx context.Context, userID string, completed bool) error
}

//go:generate mockery --name TokenRepository --output ../../mocks --outpkg mocks --case underscore
type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, rt *RefreshToken) error
	GetRefreshToken(ctx context.Context, id string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id string) error
}

//go:generate mockery --name LibraryRepository --output ../../mocks --outpkg mocks --case underscore
type LibraryRepository interface {
	AddLikedBook(ctx context.Context, userID, bookID string) error
	RemoveLikedBook(ctx context.Context, userID, bookID string) error
	GetLikedBooks(ctx context.Context, userID string) ([]LikedBook, error)
	AddLikedMovie(ctx context.Context, userID, movieID string) error
	RemoveLikedMovie(ctx context.Context, userID, movieID string) error
	GetLikedMovies(ctx context.Context, userID string) ([]LikedMovie, error)
}

//go:generate mockery --name RecommendationRepository --output ../../mocks --outpkg mocks --case underscore
type RecommendationRepository interface {
	SaveHistory(ctx context.Context, entry *RecommendationHistory) error
	GetHistory(ctx context.Context, userID string) ([]RecommendationHistory, error)
	GetHistoryByTaskID(ctx context.Context, taskID string) (*RecommendationHistory, error)
	SaveRecommendation(ctx context.Context, rec *SavedRecommendation) error
	DeleteSavedRecommendation(ctx context.Context, id string) error
	GetSavedRecommendations(ctx context.Context, userID string) ([]SavedRecommendation, error)
	UpdateResult(ctx context.Context, taskID string, resultJSON string) error // ← добавить
}

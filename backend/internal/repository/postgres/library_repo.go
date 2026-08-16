package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/pkg/logger"
	"go.uber.org/zap"
)

type LibraryRepo struct {
	pool Querier
}

func NewLibraryRepo(pool Querier) domain.LibraryRepository {
	return &LibraryRepo{pool: pool}
}

// --- книги ---
func (r *LibraryRepo) AddLikedBook(ctx context.Context, userID, bookID string) error {
	log := logger.FromContext(ctx)
	log.Debug("adding liked book", zap.String("user_id", userID), zap.String("book_id", bookID))
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_liked_books (user_id, book_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, bookID,
	)
	if err != nil {
		log.Error("failed to add liked book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
	}
	return err
}

func (r *LibraryRepo) RemoveLikedBook(ctx context.Context, userID, bookID string) error {
	log := logger.FromContext(ctx)
	log.Debug("removing liked book", zap.String("user_id", userID), zap.String("book_id", bookID))
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_liked_books WHERE user_id = $1 AND book_id = $2`,
		userID, bookID,
	)
	if err != nil {
		log.Error("failed to remove liked book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
	}
	return err
}

func (r *LibraryRepo) GetLikedBooks(ctx context.Context, userID string) ([]domain.LikedBook, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting liked books", zap.String("user_id", userID))
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, book_id, liked_at FROM user_liked_books WHERE user_id = $1 ORDER BY liked_at DESC`,
		userID,
	)
	if err != nil {
		log.Error("failed to query liked books", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var books []domain.LikedBook
	for rows.Next() {
		var b domain.LikedBook
		if err := rows.Scan(&b.UserID, &b.BookID, &b.LikedAt); err != nil {
			log.Error("failed to scan liked book", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("liked books retrieved", zap.String("user_id", userID), zap.Int("count", len(books)))
	return books, nil
}

// --- фильмы ---
func (r *LibraryRepo) AddLikedMovie(ctx context.Context, userID, movieID string) error {
	log := logger.FromContext(ctx)
	log.Debug("adding liked movie", zap.String("user_id", userID), zap.String("movie_id", movieID))
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_liked_movies (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, movieID,
	)
	if err != nil {
		log.Error("failed to add liked movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
	}
	return err
}

func (r *LibraryRepo) RemoveLikedMovie(ctx context.Context, userID, movieID string) error {
	log := logger.FromContext(ctx)
	log.Debug("removing liked movie", zap.String("user_id", userID), zap.String("movie_id", movieID))
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_liked_movies WHERE user_id = $1 AND movie_id = $2`,
		userID, movieID,
	)
	if err != nil {
		log.Error("failed to remove liked movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
	}
	return err
}

func (r *LibraryRepo) GetLikedMovies(ctx context.Context, userID string) ([]domain.LikedMovie, error) {
	log := logger.FromContext(ctx)
	log.Debug("getting liked movies", zap.String("user_id", userID))
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, movie_id, liked_at FROM user_liked_movies WHERE user_id = $1 ORDER BY liked_at DESC`,
		userID,
	)
	if err != nil {
		log.Error("failed to query liked movies", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var movies []domain.LikedMovie
	for rows.Next() {
		var m domain.LikedMovie
		if err := rows.Scan(&m.UserID, &m.MovieID, &m.LikedAt); err != nil {
			log.Error("failed to scan liked movie", zap.Error(err), zap.String("user_id", userID))
			return nil, err
		}
		movies = append(movies, m)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows iteration error", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}
	log.Debug("liked movies retrieved", zap.String("user_id", userID), zap.Int("count", len(movies)))
	return movies, nil
}

package postgres

import (
	"context"

	"github.com/I000000/recly/internal/domain"
)

type LibraryRepo struct {
	pool Querier
}

func NewLibraryRepo(pool Querier) domain.LibraryRepository {
	return &LibraryRepo{pool: pool}
}

// --- книги ---
func (r *LibraryRepo) AddLikedBook(ctx context.Context, userID, bookID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_liked_books (user_id, book_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, bookID,
	)
	return err
}

func (r *LibraryRepo) RemoveLikedBook(ctx context.Context, userID, bookID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_liked_books WHERE user_id = $1 AND book_id = $2`,
		userID, bookID,
	)
	return err
}

func (r *LibraryRepo) GetLikedBooks(ctx context.Context, userID string) ([]domain.LikedBook, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, book_id, liked_at FROM user_liked_books WHERE user_id = $1 ORDER BY liked_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.LikedBook
	for rows.Next() {
		var b domain.LikedBook
		if err := rows.Scan(&b.UserID, &b.BookID, &b.LikedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// --- фильмы ---
func (r *LibraryRepo) AddLikedMovie(ctx context.Context, userID, movieID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_liked_movies (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, movieID,
	)
	return err
}

func (r *LibraryRepo) RemoveLikedMovie(ctx context.Context, userID, movieID string) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_liked_movies WHERE user_id = $1 AND movie_id = $2`,
		userID, movieID,
	)
	return err
}

func (r *LibraryRepo) GetLikedMovies(ctx context.Context, userID string) ([]domain.LikedMovie, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, movie_id, liked_at FROM user_liked_movies WHERE user_id = $1 ORDER BY liked_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []domain.LikedMovie
	for rows.Next() {
		var m domain.LikedMovie
		if err := rows.Scan(&m.UserID, &m.MovieID, &m.LikedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, rows.Err()
}

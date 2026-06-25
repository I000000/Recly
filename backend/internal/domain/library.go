package domain

import "time"

type LikedBook struct {
	UserID  string    `json:"user_id"`
	BookID  string    `json:"book_id"`
	LikedAt time.Time `json:"liked_at"`
}

type LikedMovie struct {
	UserID  string    `json:"user_id"`
	MovieID string    `json:"movie_id"`
	LikedAt time.Time `json:"liked_at"`
}

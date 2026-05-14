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

type RecommendationHistory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TaskID      string    `json:"task_id"`
	SelectedIDs []string  `json:"selected_ids"`
	Direction   string    `json:"direction"`
	Weights     string    `json:"weights"`
	Result      string    `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
}

type SavedRecommendation struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	FromType string    `json:"from_type"`
	FromID   string    `json:"from_id"`
	ToType   string    `json:"to_type"`
	ToID     string    `json:"to_id"`
	SavedAt  time.Time `json:"saved_at"`
}

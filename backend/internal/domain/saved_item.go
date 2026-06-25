package domain

import "time"

type SavedItem struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	ItemType string    `json:"item_type"`
	ItemID   string    `json:"item_id"`
	SavedAt  time.Time `json:"saved_at"`
}

package domain

import "time"

type ViewedItem struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	ItemType string    `json:"item_type"`
	ItemID   string    `json:"item_id"`
	ViewedAt time.Time `json:"viewed_at"`
}

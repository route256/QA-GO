package models

import "time"

// Item is a struct
type Item struct {
	ID        string     `json:"id"`
	Platform  string     `json:"platform"`
	UserID    string     `json:"userId"`
	EnteredAt *time.Time `json:"enteredAt"`
}

// ItemRequest is a struct
type ItemRequest struct {
	Platform string `json:"platform"`
	UserID   string `json:"userId"`
}

// ListItemsResponse is a struct
type ListItemsResponse struct {
	Items []struct {
		Item
	} `json:"items"`
}

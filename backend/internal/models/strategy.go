package models

import "time"

type Strategy struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"` // Optional field
	CreatedAt   time.Time `json:"created_at"`
}

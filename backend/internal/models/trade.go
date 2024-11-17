package models

import "time"

type Trade struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	TradeDate time.Time `json:"trade_date"`
	Strategy  string    `json:"strategy,omitempty"`
	Note      string    `json:"note,omitempty"`
}

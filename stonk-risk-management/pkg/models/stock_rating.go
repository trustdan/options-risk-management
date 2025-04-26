package models

import (
	"time"
)

// StockRating represents a rating for a specific stock on a specific date
type StockRating struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"date"`
	Symbol         string    `json:"symbol"`         // Stock ticker symbol
	Sector         string    `json:"sector"`         // Industry sector
	StockSentiment int       `json:"stockSentiment"` // 1-10 scale
	PriceTarget    float64   `json:"priceTarget"`    // Optional price target
	Confidence     int       `json:"confidence"`     // 1-10 scale
	Notes          string    `json:"notes"`          // Optional notes about the rating
}

// NewStockRating creates a new stock rating with the current date
func NewStockRating() *StockRating {
	return &StockRating{
		Date: time.Now(),
	}
}

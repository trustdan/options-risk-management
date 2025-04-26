package models

import (
	"time"
)

// Trade represents an options trade transaction
type Trade struct {
	ID             string    `json:"id"`             // Unique identifier for the trade group (if multi-leg)
	Symbol         string    `json:"symbol"`         // Stock ticker symbol
	Sector         string    `json:"sector"`         // Industry sector
	Strategy       string    `json:"strategy"`       // Strategy category (e.g., "Vertical Spreads")
	Type           string    `json:"type"`           // Specific strategy type (e.g., "Bull Call Spread")
	Week           int       `json:"week"`           // Calendar week number for display
	EntryDate      time.Time `json:"entryDate"`      // Date the trade was entered
	ExpirationDate time.Time `json:"expirationDate"` // Date the option expires
	EntryPrice     float64   `json:"entryPrice"`     // Price paid/received for the trade
	Notes          string    `json:"notes"`          // Optional notes
}

package database

import (
	"encoding/json"
	"fmt"
	"sort"

	"stonk-risk-management/pkg/models"

	"github.com/google/uuid"
)

const tradePrefix = "trade:"

// TradeRepository handles database operations for trades
type TradeRepository struct {
	db *DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db *DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// Save saves a trade to the database
// If the trade ID already exists, it updates the existing record.
func (r *TradeRepository) Save(trade *models.Trade) error {
	// Use the existing ID if provided (for updates), otherwise generate a new one.
	if trade.ID == "" {
		trade.ID = uuid.New().String()
	}

	// Create a unique key for this trade
	// Format: trade:<ID>
	key := fmt.Sprintf("%s%s", tradePrefix, trade.ID)

	return r.db.Put(key, trade)
}

// Delete removes a trade from the database
func (r *TradeRepository) Delete(id string) error {
	key := fmt.Sprintf("%s%s", tradePrefix, id)
	return r.db.Delete(key)
}

// GetAll retrieves all trades from the database
func (r *TradeRepository) GetAll() ([]*models.Trade, error) {
	values, err := r.db.GetAllWithPrefix(tradePrefix)
	if err != nil {
		return nil, err
	}

	trades := make([]*models.Trade, 0, len(values))
	for _, v := range values {
		trade := &models.Trade{}
		if err := json.Unmarshal(v, trade); err != nil {
			// Log the problematic entry but try to continue
			fmt.Printf("Error unmarshalling trade data: %v\nData: %s\n", err, string(v))
			continue
		}
		trades = append(trades, trade)
	}

	// Sort by entry date (most recent first)
	sort.Slice(trades, func(i, j int) bool {
		return trades[i].EntryDate.After(trades[j].EntryDate)
	})

	return trades, nil
}

package database

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

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
	// Note: For multi-leg trades saved together, they share the same ID.
	// This simple Save won't group them; the frontend logic handles creating multiple entries.
	if trade.ID == "" {
		trade.ID = uuid.New().String()
	}
	// Key needs to be unique per *individual* leg, even if they share a conceptual ID.
	// We can use the ID + expiration date to make it unique per leg.
	key := fmt.Sprintf("%s%s_%s", tradePrefix, trade.ID, trade.ExpirationDate.Format(time.RFC3339Nano))
	return r.db.Put(key, trade)
}

// Delete removes a trade (or all legs of a trade) from the database
// Since legs are saved under unique keys (ID + Expiration), we need to find all keys for a given trade ID.
func (r *TradeRepository) Delete(id string) error {
	keysToDelete, err := r.db.GetKeysWithPrefix(fmt.Sprintf("%s%s_", tradePrefix, id))
	if err != nil {
		return fmt.Errorf("failed to find keys for trade ID %s: %w", id, err)
	}

	if len(keysToDelete) == 0 {
		// No keys found for this ID, maybe already deleted or invalid ID
		return nil // Or return an error like ErrNotFound
	}

	for _, key := range keysToDelete {
		if err := r.db.Delete(key); err != nil {
			// Log the error but try to continue deleting others
			fmt.Printf("Error deleting key %s for trade ID %s: %v\n", key, id, err)
		}
	}
	return nil
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

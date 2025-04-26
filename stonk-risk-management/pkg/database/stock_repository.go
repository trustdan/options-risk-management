package database

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"stonk-risk-management/pkg/models"

	"github.com/google/uuid"
)

const stockPrefix = "stock:"

// StockRepository handles database operations for stock ratings
type StockRepository struct {
	db *DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *DB) *StockRepository {
	return &StockRepository{db: db}
}

// Save saves a stock rating to the database
func (r *StockRepository) Save(rating *models.StockRating) error {
	if rating.ID == "" {
		rating.ID = uuid.New().String()
	}
	key := fmt.Sprintf("%s%s", stockPrefix, rating.ID)
	return r.db.Put(key, rating)
}

// Get retrieves a stock rating by ID
func (r *StockRepository) Get(id string) (*models.StockRating, error) {
	key := fmt.Sprintf("%s%s", stockPrefix, id)
	rating := &models.StockRating{}
	err := r.db.Get(key, rating)
	if err != nil {
		return nil, err
	}
	return rating, nil
}

// Delete removes a stock rating from the database
func (r *StockRepository) Delete(id string) error {
	key := fmt.Sprintf("%s%s", stockPrefix, id)
	return r.db.Delete(key)
}

// GetAll retrieves all stock ratings
func (r *StockRepository) GetAll() ([]*models.StockRating, error) {
	values, err := r.db.GetAllWithPrefix(stockPrefix)
	if err != nil {
		return nil, err
	}

	ratings := make([]*models.StockRating, 0, len(values))
	for _, v := range values {
		rating := &models.StockRating{}
		if err := json.Unmarshal(v, rating); err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	// Sort by date
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Date.Before(ratings[j].Date)
	})

	return ratings, nil
}

// GetByDate retrieves stock ratings for a specific date
func (r *StockRepository) GetByDate(date time.Time) ([]*models.StockRating, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	// Format date to compare only year, month, day
	targetDate := date.Format("2006-01-02")

	var filtered []*models.StockRating
	for _, a := range all {
		if a.Date.Format("2006-01-02") == targetDate {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

// GetBySector retrieves stock ratings for a specific sector
func (r *StockRepository) GetBySector(sector string) ([]*models.StockRating, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []*models.StockRating
	for _, a := range all {
		if a.Sector == sector {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

// GetBySymbol retrieves stock ratings for a specific symbol
func (r *StockRepository) GetBySymbol(symbol string) ([]*models.StockRating, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []*models.StockRating
	for _, a := range all {
		if a.Symbol == symbol {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"stonk-risk-management/pkg/database"
	"stonk-risk-management/pkg/models"

	"github.com/dgraph-io/badger/v3"
)

// App struct
type App struct {
	ctx             context.Context
	db              *database.DB
	riskRepository  *database.RiskRepository
	stockRepository *database.StockRepository
	tradeRepository *database.TradeRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Setup database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, ".options-risk-management")
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	a.db = db
	a.riskRepository = database.NewRiskRepository(db)
	a.stockRepository = database.NewStockRepository(db)
	a.tradeRepository = database.NewTradeRepository(db)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

// GetRiskAssessments returns all risk assessments
func (a *App) GetRiskAssessments() ([]*models.RiskAssessment, error) {
	return a.riskRepository.GetAll()
}

// SaveRiskAssessment saves a risk assessment
func (a *App) SaveRiskAssessment(assessment *models.RiskAssessment) error {
	return a.riskRepository.Save(assessment)
}

// DeleteRiskAssessment deletes a risk assessment
func (a *App) DeleteRiskAssessment(id string) error {
	return a.riskRepository.Delete(id)
}

// GetStockRatings returns all stock ratings
func (a *App) GetStockRatings() ([]*models.StockRating, error) {
	return a.stockRepository.GetAll()
}

// GetStockRatingsByDate returns stock ratings for a specific date
func (a *App) GetStockRatingsByDate(date time.Time) ([]*models.StockRating, error) {
	return a.stockRepository.GetByDate(date)
}

// SaveStockRating saves a stock rating
func (a *App) SaveStockRating(rating *models.StockRating) error {
	return a.stockRepository.Save(rating)
}

// DeleteStockRating deletes a stock rating
func (a *App) DeleteStockRating(id string) error {
	return a.stockRepository.Delete(id)
}

// GetTrades returns all trades
func (a *App) GetTrades() ([]*models.Trade, error) {
	return a.tradeRepository.GetAll()
}

// SaveTrade saves a trade
func (a *App) SaveTrade(trade *models.Trade) error {
	// Basic validation before saving
	if trade.Symbol == "" || trade.Sector == "" || trade.Strategy == "" || trade.Type == "" {
		return fmt.Errorf("invalid trade data: missing required fields")
	}

	// For backward compatibility, always set legNumber to 1
	trade.LegNumber = 1

	return a.tradeRepository.Save(trade)
}

// DeleteTrade deletes all legs associated with a trade ID
func (a *App) DeleteTrade(id string) error {
	return a.tradeRepository.Delete(id)
}

// RunDatabaseMaintenance performs database maintenance tasks including garbage collection
// Returns a success message or error message
func (a *App) RunDatabaseMaintenance() string {
	if a.db == nil {
		return "Error: Database not initialized"
	}

	err := a.db.RunGC()
	if err != nil {
		if err == badger.ErrNoRewrite {
			return "No garbage collection needed at this time. Database is already optimized."
		}
		return fmt.Sprintf("Error during garbage collection: %v", err)
	}

	return "Database maintenance completed successfully. Freed up unused space."
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

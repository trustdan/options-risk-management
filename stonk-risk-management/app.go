package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"stonk-risk-management/pkg/database"
	"stonk-risk-management/pkg/models"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileFilter defines a filter for file dialogs
type FileFilter struct {
	DisplayName string   `json:"name"`
	Pattern     []string `json:"extensions"`
}

// DialogOptions contains common options for file dialogs
type DialogOptions struct {
	Title           string       `json:"title"`
	DefaultFilename string       `json:"defaultFilename"`
	Filters         []FileFilter `json:"filters"`
}

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

// SaveTrade saves a single trade leg
// Note: Frontend is responsible for creating/updating all necessary legs for multi-leg trades.
func (a *App) SaveTrade(trade *models.Trade) error {
	// Basic validation before saving
	if trade.Symbol == "" || trade.Sector == "" || trade.Strategy == "" || trade.Type == "" {
		return fmt.Errorf("invalid trade data: missing required fields")
	}
	return a.tradeRepository.Save(trade)
}

// DeleteTrade deletes all legs associated with a trade ID
func (a *App) DeleteTrade(id string) error {
	return a.tradeRepository.Delete(id)
}

// GetCurrentDirectory returns the current working directory
func (a *App) GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetUserDirectory returns the user's home directory
func (a *App) GetUserDirectory() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// ShowSaveDialog shows a save file dialog using the Wails runtime
func (a *App) ShowSaveDialog(options map[string]interface{}) (string, error) {
	if a.ctx == nil {
		return "", errors.New("context is nil, app not initialized properly")
	}

	title, _ := options["title"].(string)
	defaultFilename, _ := options["defaultFilename"].(string)

	// Convert filters to the expected format
	var filters []runtime.FileFilter
	if filtersRaw, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filtersRaw {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				name, _ := filterMap["name"].(string)

				var extensions []string
				if exts, ok := filterMap["extensions"].([]interface{}); ok {
					for _, ext := range exts {
						if extStr, ok := ext.(string); ok {
							extensions = append(extensions, extStr)
						}
					}
				}

				filters = append(filters, runtime.FileFilter{
					DisplayName: name,
					Pattern:     extensions,
				})
			}
		}
	}

	// Show save dialog
	result, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters:         filters,
	})

	return result, err
}

// ShowOpenDialog shows an open file dialog using the Wails runtime
func (a *App) ShowOpenDialog(options map[string]interface{}) (string, error) {
	if a.ctx == nil {
		return "", errors.New("context is nil, app not initialized properly")
	}

	title, _ := options["title"].(string)

	// Convert filters to the expected format
	var filters []runtime.FileFilter
	if filtersRaw, ok := options["filters"].([]interface{}); ok {
		for _, filter := range filtersRaw {
			if filterMap, ok := filter.(map[string]interface{}); ok {
				name, _ := filterMap["name"].(string)

				var extensions []string
				if exts, ok := filterMap["extensions"].([]interface{}); ok {
					for _, ext := range exts {
						if extStr, ok := ext.(string); ok {
							extensions = append(extensions, extStr)
						}
					}
				}

				filters = append(filters, runtime.FileFilter{
					DisplayName: name,
					Pattern:     extensions,
				})
			}
		}
	}

	// Show open dialog
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: filters,
	})

	return result, err
}

// ExportData exports all application data to a JSON file
func (a *App) ExportData(filePath string) error {
	// If the path is not absolute, make it relative to the current directory
	if !filepath.IsAbs(filePath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		filePath = filepath.Join(currentDir, filePath)
	}

	// For now, we'll create a dummy data structure to export
	exportData := map[string]interface{}{
		"version":    "1.0",
		"exportDate": time.Now().Format(time.RFC3339),
		"preferences": map[string]interface{}{
			"theme": "dark",
			// Add more preferences here
		},
		"positions": []map[string]interface{}{
			{
				"id":           "1",
				"symbol":       "AAPL",
				"strategy":     "Bull Call Spread",
				"openDate":     "2023-04-15",
				"closeDate":    "2023-04-30",
				"profitLoss":   1200.50,
				"commissions":  12.99,
				"notes":        "Closed early due to earnings announcement",
				"legs":         []interface{}{},
				"adjustments":  []interface{}{},
				"screenshots":  []interface{}{},
				"riskAnalysis": map[string]interface{}{},
			},
			// Add more positions here
		},
		"tradeJournal": []map[string]interface{}{
			{
				"id":        "j1",
				"date":      "2023-04-15",
				"title":     "Market Analysis",
				"content":   "Today I observed the following patterns...",
				"tags":      []string{"SPY", "analysis", "patterns"},
				"sentiment": "bullish",
			},
			// Add more journal entries here
		},
	}

	// Convert the data to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	// Write to file
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// ImportData imports application data from a JSON file
func (a *App) ImportData(filePath string) error {
	// If the path is not absolute, make it relative to the current directory
	if !filepath.IsAbs(filePath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		filePath = filepath.Join(currentDir, filePath)
	}

	// Read file content
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse JSON
	var importedData map[string]interface{}
	err = json.Unmarshal(data, &importedData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Here you would process the imported data and update your application state
	// For now, we'll just return success
	return nil
}

// Load the app menu during app startup
func (a *App) getAppMenu() *menu.Menu {
	appMenu := menu.NewMenu()

	// The File Menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Import Data", keys.CmdOrCtrl("I"), func(_ *menu.CallbackData) {
		// TODO: Implement Import
	})
	fileMenu.AddText("Export Data", keys.CmdOrCtrl("E"), func(_ *menu.CallbackData) {
		// TODO: Implement Export
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Quit", keys.CmdOrCtrl("Q"), func(_ *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})

	// The Edit Menu
	editMenu := appMenu.AddSubmenu("Edit")
	editMenu.AddText("Cut", keys.CmdOrCtrl("X"), nil)
	editMenu.AddText("Copy", keys.CmdOrCtrl("C"), nil)
	editMenu.AddText("Paste", keys.CmdOrCtrl("V"), nil)
	editMenu.AddSeparator()
	editMenu.AddText("Preferences", keys.CmdOrCtrl("P"), func(_ *menu.CallbackData) {
		// TODO: Show preferences dialog
	})

	// The Help Menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("About", nil, func(_ *menu.CallbackData) {
		// TODO: Show about dialog
	})
	helpMenu.AddText("Documentation", keys.CmdOrCtrl("H"), nil)

	return appMenu
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

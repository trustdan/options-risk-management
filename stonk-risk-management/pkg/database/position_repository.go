package database

import (
	"stonk-risk-management/pkg/models"

	"github.com/dgraph-io/badger/v3"
)

const positionSettingsKey = "position_settings"

// PositionRepository handles database operations for position settings
type PositionRepository struct {
	db *DB
}

// NewPositionRepository creates a new position repository
func NewPositionRepository(db *DB) *PositionRepository {
	return &PositionRepository{db: db}
}

// GetSettings retrieves the user's position settings
func (r *PositionRepository) GetSettings() (*models.PositionSettings, error) {
	settings := &models.PositionSettings{}
	err := r.db.Get(positionSettingsKey, settings)

	if err != nil {
		if err == badger.ErrKeyNotFound {
			// Return default settings if none exist yet
			return &models.PositionSettings{
				AccountValue:          25000,
				AccountRiskPerTrade:   2,
				MaxPortfolioExposure:  20,
				StopLossPercent:       50,
				RiskRewardRatio:       1,
				DailyLossLimit:        3,
				WeeklyLossLimit:       7,
				PositionScaling:       100,
				CorrelationAdjustment: 1,
				VolatilityMultiplier:  1,
				MaxDrawdownTolerance:  15,
			}, nil
		}
		return nil, err
	}

	return settings, nil
}

// SaveSettings saves the user's position settings
func (r *PositionRepository) SaveSettings(settings *models.PositionSettings) error {
	return r.db.Put(positionSettingsKey, settings)
}

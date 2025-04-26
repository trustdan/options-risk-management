package database

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"stonk-risk-management/pkg/models"

	"github.com/google/uuid"
)

const riskPrefix = "risk:"

// RiskRepository handles database operations for risk assessments
type RiskRepository struct {
	db *DB
}

// NewRiskRepository creates a new risk repository
func NewRiskRepository(db *DB) *RiskRepository {
	return &RiskRepository{db: db}
}

// Save saves a risk assessment to the database
func (r *RiskRepository) Save(assessment *models.RiskAssessment) error {
	if assessment.ID == "" {
		assessment.ID = uuid.New().String()
	}
	key := fmt.Sprintf("%s%s", riskPrefix, assessment.ID)
	return r.db.Put(key, assessment)
}

// Get retrieves a risk assessment by ID
func (r *RiskRepository) Get(id string) (*models.RiskAssessment, error) {
	key := fmt.Sprintf("%s%s", riskPrefix, id)
	assessment := &models.RiskAssessment{}
	err := r.db.Get(key, assessment)
	if err != nil {
		return nil, err
	}
	return assessment, nil
}

// Delete removes a risk assessment from the database
func (r *RiskRepository) Delete(id string) error {
	key := fmt.Sprintf("%s%s", riskPrefix, id)
	return r.db.Delete(key)
}

// GetAll retrieves all risk assessments
func (r *RiskRepository) GetAll() ([]*models.RiskAssessment, error) {
	values, err := r.db.GetAllWithPrefix(riskPrefix)
	if err != nil {
		return nil, err
	}

	assessments := make([]*models.RiskAssessment, 0, len(values))
	for _, v := range values {
		assessment := &models.RiskAssessment{}
		if err := json.Unmarshal(v, assessment); err != nil {
			return nil, err
		}
		assessments = append(assessments, assessment)
	}

	// Sort by date with newest last
	sort.Slice(assessments, func(i, j int) bool {
		return assessments[i].Date.Before(assessments[j].Date)
	})

	return assessments, nil
}

// GetByDateRange retrieves risk assessments within a date range
func (r *RiskRepository) GetByDateRange(start, end time.Time) ([]*models.RiskAssessment, error) {
	all, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []*models.RiskAssessment
	for _, a := range all {
		if (a.Date.Equal(start) || a.Date.After(start)) &&
			(a.Date.Equal(end) || a.Date.Before(end)) {
			filtered = append(filtered, a)
		}
	}

	return filtered, nil
}

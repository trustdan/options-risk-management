package models

import (
	"time"
)

// RiskAssessment represents a trader's daily risk assessment
type RiskAssessment struct {
	ID             string    `json:"id"`
	Date           time.Time `json:"date"`
	EmotionalScore int       `json:"emotionalScore"` // 1-10 scale
	FOMOScore      int       `json:"fomoScore"`      // 1-10 scale
	BiasScore      int       `json:"biasScore"`      // 1-10 scale
	OverallScore   int       `json:"overallScore"`   // Computed or manually entered
	Notes          string    `json:"notes"`          // Optional trader notes
}

// NewRiskAssessment creates a new risk assessment with the current date
func NewRiskAssessment() *RiskAssessment {
	return &RiskAssessment{
		Date: time.Now(),
	}
}

// CalculateOverall computes the overall risk score based on individual components
func (r *RiskAssessment) CalculateOverall() {
	// Simple average calculation, can be replaced with more complex weighting
	r.OverallScore = (r.EmotionalScore + r.FOMOScore + r.BiasScore) / 3
}

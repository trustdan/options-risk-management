package models

// PositionSettings represents the user's risk management settings for position sizing
type PositionSettings struct {
	AccountValue          float64 `json:"accountValue"`
	AccountRiskPerTrade   float64 `json:"accountRiskPerTrade"`
	MaxPortfolioExposure  float64 `json:"maxPortfolioExposure"`
	StopLossPercent       float64 `json:"stopLossPercent"`
	RiskRewardRatio       float64 `json:"riskRewardRatio"`
	DailyLossLimit        float64 `json:"dailyLossLimit"`
	WeeklyLossLimit       float64 `json:"weeklyLossLimit"`
	PositionScaling       float64 `json:"positionScaling"`
	CorrelationAdjustment float64 `json:"correlationAdjustment"`
	VolatilityMultiplier  float64 `json:"volatilityMultiplier"`
	MaxDrawdownTolerance  float64 `json:"maxDrawdownTolerance"`
}

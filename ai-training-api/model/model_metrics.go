package model

import (
	"gorm.io/gorm"
)

type ModelMetrics struct {
	// Unique identifier for the metric entry.
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	// Tenant ID is used to identify the tenant to which the metric belongs.
	Tenant       string `json:"tenant" gorm:"size:255;not null"`
	// Run ID is used to identify the specific run.
	Run          string `json:"run" gorm:"size:255;not null"`
	// The name of the metric.
	MetricName   string `json:"metric_name" gorm:"size:255;not null"`
	// Step represents the measurement step; unsigned integer.
	Step         uint   `json:"step" gorm:"not null"`
	// MetricValue stores the value of the metric as a string.
	MetricValue  string `json:"metric_value" gorm:"size:255;not null"`
}

// Add a custom hook if necessary for additional logic.
// Example: AfterCreate hook for custom logic
func (m *ModelMetrics) AfterCreate(tx *gorm.DB) error {
	// Custom logic after creating a metric entry
	tx.Logger.Info(tx.Statement.Context, "AfterCreate hook called for ModelMetrics")
	return nil
}

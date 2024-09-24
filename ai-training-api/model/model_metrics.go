package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelMetrics struct {
    StackID     uint64   `json:"stack_id" gorm:"not null;primaryKey"`
    ProcessID   uuid.UUID `json:"process_id" gorm:"type:char(36);not null;primaryKey;foreignKey:ProcessID;references:ID"` // Foreign key
    MetricName  string   `json:"metric_name" gorm:"size:32;not null;primaryKey"`
    StepName    string   `json:"step_name" gorm:"size:32;not null;primaryKey"`
    Step        uint32   `json:"step" gorm:"not null;primaryKey"`
    MetricValue string   `json:"metric_value" gorm:"size:64;not null"`

    Process Process `gorm:"foreignKey:ProcessID;references:ID"` // Relationship definition
}
// Add a custom hook if necessary for additional logic.
// Example: AfterCreate hook for custom logic
func (m *ModelMetrics) AfterCreate(tx *gorm.DB) error {
	// Custom logic after creating a metric entry
	tx.Logger.Info(tx.Statement.Context, "AfterCreate hook called for ModelMetrics")
	return nil
}

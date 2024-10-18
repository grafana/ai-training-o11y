package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StackID is what user we are using
// ProcessID is the process sending metrics by uuid
// MetricName is what metric we are logging (e.g., accuracy, loss)
// StepName is the name of the step in the process (e.g., step, batch, epoch)
// Separating out by these is important because it only makes sense to graph
// data for the same metric and step in one panel
// Step is the step number, which goes on the x-axis, and MetricValue is the y-value.
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

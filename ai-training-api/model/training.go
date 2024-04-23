package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// The database model used to track training information.
type Training struct {
	// UUID generated for the training.
	ID uuid.UUID `json:"id" gorm:"primarykey;type:char(36)" validate:"isdefault"`
	// Tenant ID is used to identify the tenant to which the training belongs.
	TenantID string `json:"tenant_id"`
	// The training name.
	Name string `json:"name"`
	// Description for the training.
	Description string `json:"description"`
	// The training status.
	Status string `json:"status"`
	// Start time.
	StartTime time.Time `json:"start_time"`
	// End time.
	EndTime time.Time `json:"end_time"`

	// A training can have multiple processes.
	Processes []Process `json:"processes" gorm:"serializer:json"`

	// Process Metadata.
	// links to HF, GH repos, DVC checkpoints, hyperparameters - can be in dozens, rarely hundreds
	// TODO: cap at 1024?
	Metadata datatypes.JSON `json:"metadata"`
}

// Uodate StartTime and EndTime based on Start and End time of the processes.
func (t *Training) UpdateTimes() {
	if len(t.Processes) == 0 {
		return
	}

	t.StartTime = t.Processes[0].StartTime
	t.EndTime = t.Processes[0].EndTime

	for _, p := range t.Processes {
		if p.StartTime.Before(t.StartTime) {
			t.StartTime = p.StartTime
		}
		if p.EndTime.After(t.EndTime) {
			t.EndTime = p.EndTime
		}
	}
}

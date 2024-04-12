package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// The database model used to track training information.
type Training struct {
	// Tenant ID is used to identify the tenant to which the training belongs.
	TenantID string `json:"tenant_id"`
	// UUID generated for the training.
	ID uuid.UUID `json:"pid"`
	// The training status.
	Status string `json:"status"`

	// A training can have multiple processes.
	Processes []Process `json:"processes" gorm:"serializer:json"`

	// Process Metadata.
	// links to HF, GH repos, DVC checkpoints, hyperparameters - can be in dozens, rarely hundreds
	// TODO: cap at 1024?
	Metadata datatypes.JSON `json:"metadata"`
}

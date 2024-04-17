package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Process struct {
	// Tenant ID is used to identify the tenant to which the process belongs.
	TenantID string `json:"tenant_id"`
	// UUID generated for the process.
	ID uuid.UUID `json:"process_uuid"`
	// The process status.
	Status string `json:"status"`

	// Process Metadata.
	// This field is used to store additional metadata about the process.
	// TODO: cap at 1024?
	Metadata datatypes.JSON `json:"metadata"`
}

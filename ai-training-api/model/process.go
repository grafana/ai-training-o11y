package model

import (
	"time"

	"github.com/google/uuid"
)

type Process struct {
	// UUID generated for the process.
	ID uuid.UUID `json:"process_uuid" gorm:"primarykey;type:char(36)" validate:"isdefault"`
	// Tenant ID is used to identify the tenant to which the process belongs.
	TenantID string `json:"tenant_id"`
	// The process status.
	Status string `json:"status"`
	// Start time.
	StartTime time.Time `json:"start_time"`
	// End time.
	EndTime time.Time `json:"end_time"`

	// Group ID is the UUID of the group to which the process belongs.
	// Its the foreign key to the Group table. It is a pointer to allow for null values.
	GroupID *uuid.UUID `json:"group_uuid" gorm:"type:char(36)"`

	Project string `json:"project"`

	// Process Metadata.
	// This field is used to store additional metadata about the process.
	// TODO: cap at 1024?
	// We are storing this in a separate table, so not serializing it here.
	Metadata []MetadataKV `json:"metadata" gorm:"-"`
}

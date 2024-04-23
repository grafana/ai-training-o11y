package model

import "github.com/google/uuid"

// MetadataKV is the database model used to track metadata information.
// This is used to flatten JSON metadata into a key-value pair and index
// it for search.
type MetadataKV struct {
	// Tenant ID is used to identify the tenant to which the metadata belongs.
	TenantID string `json:"tenant_id"`
	// Key is the metadata key.
	Key string `json:"key"`
	// Value is the metadata value.
	Value any `json:"value"`

	// Process ID is the UUID of the process to which the metadata belongs.
	// Its the foreign key to the Process table.
	ProcessID uuid.UUID `json:"process_id"`
}

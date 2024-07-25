package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// The database model used to track Group information.
type Group struct {
	// UUID generated for the group.
	ID uuid.UUID `json:"id" gorm:"primarykey;type:char(36)" validate:"isdefault"`
	// Tenant ID is used to identify the tenant to which the group belongs.
	TenantID string `json:"tenant_id"`
	// The group name.
	Name string `json:"name"`
	// Description for the group.
	Description string `json:"description"`
	// The status of processes in the group.
	Status string `json:"status"`
	// Start time.
	StartTime time.Time `json:"start_time"`
	// End time. Should be nullable to allow for groups that are still running.
	EndTime sql.NullTime `json:"end_time"`

	// Processes in the group.
	Processes []Process `json:"processes" gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// Metadata for the group.
	// links to HF, GH repos, DVC checkpoints, hyperparameters - can be in dozens, rarely hundreds
	// TODO: cap at 1024?
	Metadata datatypes.JSON `json:"metadata"`
}

// Uodate StartTime and EndTime based on Start and End time of the processes.
func (t *Group) UpdateTimes(processes []Process) {
	if len(processes) == 0 {
		return
	}

	t.StartTime = processes[0].StartTime
	t.EndTime = processes[0].EndTime

	for _, p := range processes {
		if p.StartTime.Before(t.StartTime) {
			t.StartTime = p.StartTime
		}
		if p.EndTime.Time.After(t.EndTime.Time) {
			t.EndTime = p.EndTime
		}
	}
}

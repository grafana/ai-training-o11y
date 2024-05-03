package model

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
)

// MetadataKV is the database model used to track metadata information.
// This is used to flatten JSON metadata into a key-value pair and index
// it for search.
type MetadataKV struct {
	// Tenant ID is used to identify the tenant to which the metadata belongs.
	TenantID string `json:"tenant_id"`
	// Key is the metadata key.
	Key string `json:"key"`
	// Value is the metadata value.
	Value []byte `json:"value"`
	// Type is the type of the metadata value.
	Type string `json:"type"`

	// Process ID is the UUID of the process to which the metadata belongs.
	// Its the foreign key to the Process table.
	ProcessID uuid.UUID `json:"process_id"`
}

func MarshalMetadataValue(value interface{}) (string, []byte) {
	switch v := value.(type) {
	case string:
		return "string", []byte(v)
	case int:
		return "int", []byte(fmt.Sprintf("%d", v))
	case float64:
		// https://stackoverflow.com/a/55436758
		if value.(float64) == float64(int64(value.(float64))) {
			return "int", []byte(fmt.Sprintf("%d", int64(value.(float64))))
		}

		bits := math.Float64bits(v)
		byteArr := make([]byte, 8)
		binary.BigEndian.PutUint64(byteArr, bits)
		return "float", byteArr
	case bool:
		return "bool", []byte(strconv.FormatBool(v))
	default:
		return "unknown", nil
	}
}

func UnmarshalMetadataValue(value []byte, valueType string) (interface{}, error) {
	switch valueType {
	case "string":
		return string(value), nil
	case "int":
		return strconv.Atoi(string(value))
	case "float":
		bits := binary.BigEndian.Uint64(value)
		float := math.Float64frombits(bits)
		return float, nil
	case "bool":
		return strconv.ParseBool(string(value))
	default:
		return nil, fmt.Errorf("unknown type: %s", valueType)
	}
}

package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB map[string]any

// Scan implementasi untuk membaca dari database (jsonb -> map)
func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = make(map[string]any)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for JSONB: %T", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementasi untuk menyimpan ke database (map -> jsonb)
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

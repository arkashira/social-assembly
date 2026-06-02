package federation

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Activity represents a subset of ActivityPub activities we need to support
type Activity struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Actor     string          `json:"actor"`
	Object    json.RawMessage `json:"object"`
	Target    string          `json:"target,omitempty"`
	Published *time.Time      `json:"published,omitempty"`
}

// ParseActivity parses and validates an ActivityPub JSON payload
// Supported activity types: Create, Delete, Follow
// Returns parsed Activity or validation error
func ParseActivity(data []byte) (*Activity, error) {
	// First pass: validate structure and extract type
	var raw struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if raw.Type == "" {
		return nil, errors.New("missing 'type' field")
	}

	// Validate activity type
	if !isSupportedType(raw.Type) {
		return nil, fmt.Errorf("unsupported activity type: %s", raw.Type)
	}

	// Second pass: full unmarshal and validation
	var act Activity
	if err := json.Unmarshal(data, &act); err != nil {
		return nil, fmt.Errorf("invalid activity structure: %w", err)
	}

	// Validate required fields
	if err := validateRequiredFields(&act); err != nil {
		return nil, err
	}

	// Type-specific validation
	if err := validateTypeSpecificFields(&act); err != nil {
		return nil, err
	}

	return &act, nil
}

func isSupportedType(t string) bool {
	switch t {
	case "Create", "Delete", "Follow":
		return true
	default:
		return false
	}
}

func validateRequiredFields(act *Activity) error {
	if act.ID == "" {
		return errors.New("missing 'id' field")
	}
	if act.Actor == "" {
		return errors.New("missing 'actor' field")
	}
	if len(act.Object) == 0 {
		return errors.New("missing 'object' field")
	}
	return nil
}

func validateTypeSpecificFields(act *Activity) error {
	switch act.Type {
	case "Follow":
		if act.Target == "" {
			return errors.New("missing 'target' field for Follow activity")
		}
	}
	return nil
}
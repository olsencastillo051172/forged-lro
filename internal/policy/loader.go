package policy

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadPolicy reads the rotation policy from a file and deserializes it.
// It performs basic syntax and structural checks during the JSON unmarshaling process.
func LoadPolicy(path string) (*RotationPolicy, error) {
	// 1. Physical Read: Ensure the file is accessible
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("AUDIT_FAIL: could not read policy file at %s: %w", path, err)
	}

	// 2. Deserialization: Map JSON to our strictly typed structs
	var pol RotationPolicy
	if err := json.Unmarshal(data, &pol); err != nil {
		return nil, fmt.Errorf("AUDIT_FAIL: policy file has malformed JSON structure: %w", err)
	}

	// 3. Structural Check: Ensure the policy is not an empty object
	if pol.PolicyVersion == "" {
		return nil, fmt.Errorf("AUDIT_FAIL: policy file is empty or missing mandatory version field")
	}

	return &pol, nil
}

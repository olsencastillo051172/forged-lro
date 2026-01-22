package policy

import (
	"errors"
	"fmt"
)

// ValidateInvariants enforces the technical and legal boundaries of the policy.
// It ensures that the loaded configuration strictly adheres to RVA standards.
func ValidateInvariants(p *RotationPolicy) error {
	// 1. Cryptographic Invariants
	if p.Constraints.HashAlg != "sha256" {
		return fmt.Errorf("AUDIT_FAIL: hash_alg '%s' is not supported (required: sha256)", p.Constraints.HashAlg)
	}

	foundAlg := false
	for _, alg := range p.Constraints.AllowedHashAlgs {
		if alg == "sha256" {
			foundAlg = true
			break
		}
	}
	if !foundAlg {
		return errors.New("AUDIT_FAIL: sha256 must be present in allowed_hash_algs")
	}

	if p.Constraints.DomainSeparator != "RVA_NODE:v1" {
		return fmt.Errorf("AUDIT_FAIL: domain_separator '%s' violates protocol version (required: RVA_NODE:v1)", p.Constraints.DomainSeparator)
	}

	// 2. Merkle Tree Boundaries
	if p.Constraints.MinDepth < 1 || p.Constraints.MaxDepth > 64 {
		return fmt.Errorf("AUDIT_FAIL: invalid Merkle depth boundaries (min:1, max:64)")
	}

	// 3. Epoch & Timing Discipline
	// NOTE: We enforce the 24h production limit here. 
	// Developer overrides should be handled via environment variables, not by weakening the policy.
	if p.Epochs.IntervalSeconds < 86400 {
		return fmt.Errorf("AUDIT_FAIL: rotation interval %d is below production safety limit (86400s)", p.Epochs.IntervalSeconds)
	}

	if p.Epochs.IDFormat != "numeric_ascending" {
		return fmt.Errorf("AUDIT_FAIL: epoch_id_format '%s' is not recognized", p.Epochs.IDFormat)
	}

	// 4. Governance Rules (Cutover)
	if !p.Cutover.RequirePrevAnchor || !p.Cutover.StrictMonotonicEpoch {
		return errors.New("AUDIT_FAIL: cutover rules must enforce previous_anchor and strict_monotonicity")
	}

	return nil
}

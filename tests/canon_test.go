package tests

import (
	"testing"

	"github.com/olsencastillo051172/forged-lro/src/config"
)

// ─────────────────────────────────────────────
// Canon v1.0 Guardrail Tests
// ─────────────────────────────────────────────

func TestCanonVersionIsFrozen(t *testing.T) {
	if config.CanonVersion != "v1.0" {
		t.Fatalf("CanonVersion modified: expected v1.0, got %s", config.CanonVersion)
	}

	if config.CanonStatus != "FROZEN" {
		t.Fatalf("CanonStatus modified: expected FROZEN, got %s", config.CanonStatus)
	}
}

func TestEpochSizeInvariant(t *testing.T) {
	if config.EpochSize != 1000 {
		t.Fatalf("EpochSize invariant violated: expected 1000, got %d", config.EpochSize)
	}
}

func TestCryptoPrimitivesInvariant(t *testing.T) {
	if config.HashAlgorithm != "SHA-256" {
		t.Fatalf("Hash algorithm modified: %s", config.HashAlgorithm)
	}

	if config.SignatureAlgorithm != "Ed25519" {
		t.Fatalf("Signature algorithm modified: %s", config.SignatureAlgorithm)
	}

	if config.TimeStandard != "UTC" {
		t.Fatalf("Time standard modified: %s", config.TimeStandard)
	}
}

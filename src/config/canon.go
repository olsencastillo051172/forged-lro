package config

// FORGED-LRO — Canon v1.0
// Status: FROZEN
// Scope: Invariants only (no business logic)
// Effective date: 2026-01-10

// ─────────────────────────────────────────────
// Cryptographic primitives (NON-NEGOTIABLE)
// ─────────────────────────────────────────────

const (
	HashAlgorithm      = "SHA-256"
	SignatureAlgorithm = "Ed25519"
	TimeStandard       = "UTC"
)

// ─────────────────────────────────────────────
// Ledger & Epoch Configuration
// ─────────────────────────────────────────────

const (
	// Number of ledger entries per epoch
	EpochSize = 1000

	// Genesis rule: first entry has no previous hash
	GenesisPrevHash = ""
)

// ─────────────────────────────────────────────
// Timestamp Validation
// ─────────────────────────────────────────────

const (
	// ±5 minutes tolerance for submitted timestamps
	SubmissionTimestampToleranceSeconds = 300
)

// ─────────────────────────────────────────────
// Supply Chain Constraints (C4 Principle)
// ─────────────────────────────────────────────

const (
	// Default clock skew tolerance (seconds)
	DefaultClockSkewSeconds = 600

	// Maximum plausible speed (km/h) to prevent GPS teleportation
	MaxPhysicalSpeedKmh = 1000

	// Maximum allowed custody gap (seconds)
	CustodyGapThresholdSeconds = 21600 // 6 hours
)

// ─────────────────────────────────────────────
// Canon Metadata
// ─────────────────────────────────────────────

const (
	CanonVersion = "v1.0"
	CanonStatus  = "FROZEN"
)

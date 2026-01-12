// Package ledger provides append-only JSONL storage for FORGED-LRO Canon v1.0.
// All operations are deterministic and immutable: no UPDATE/DELETE operations.
package ledger

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var (
	// ErrLedgerIO is returned when file I/O operations fail
	ErrLedgerIO = errors.New("ledger I/O error")

	// ErrLedgerCorrupt is returned when ledger data is corrupted or invalid
	ErrLedgerCorrupt = errors.New("ledger corrupt")

	// ErrNoRegistrations is returned when attempting to seal with no registrations
	ErrNoRegistrations = errors.New("no registrations to seal")

	// ErrInvalidHex is returned when hex validation fails
	ErrInvalidHex = errors.New("invalid hex format")

	// ErrInvalidTimestamp is returned when timestamp validation fails
	ErrInvalidTimestamp = errors.New("invalid timestamp format")
)

// hex64Pattern validates 64-character lowercase hex strings (SHA-256)
var hex64Pattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

// hex128Pattern validates 128-character lowercase hex strings (Ed25519 signatures)
var hex128Pattern = regexp.MustCompile(`^[a-f0-9]{128}$`)

// Default ledger path
var (
	ledgerPath   = "data/ledger.jsonl"
	ledgerMutex  sync.Mutex
)

// SetLedgerPath allows configuration of the ledger file path (primarily for testing)
func SetLedgerPath(path string) {
	ledgerMutex.Lock()
	defer ledgerMutex.Unlock()
	ledgerPath = path
}

// GetLedgerPath returns the current ledger file path
func GetLedgerPath() string {
	ledgerMutex.Lock()
	defer ledgerMutex.Unlock()
	return ledgerPath
}

// RegisterEntry represents a registration record in the ledger
type RegisterEntry struct {
	Type             string `json:"type"`                        // Always "register"
	Canon            string `json:"canon"`                       // Canon version (e.g., "v1.0")
	Timestamp        string `json:"timestamp"`                   // RFC3339Nano format
	ObjectHashHex    string `json:"object_hash_hex"`             // 64 lowercase hex
	CanonicalJSONB64 string `json:"canonical_json_b64,omitempty"` // Optional base64 encoded canonical JSON
}

// Manifest represents the seal manifest containing cryptographic proof
type Manifest struct {
	MerkleRoot string `json:"merkle_root"` // 64 lowercase hex
	Signature  string `json:"signature"`   // 128 lowercase hex (Ed25519)
	PublicKey  string `json:"public_key"`  // 64 lowercase hex (Ed25519)
	Timestamp  string `json:"timestamp"`   // RFC3339Nano format
}

// SealEntry represents a seal record in the ledger
type SealEntry struct {
	Type     string   `json:"type"`     // Always "seal"
	Manifest Manifest `json:"manifest"`
}

// AppendRegister appends a registration entry to the ledger.
// 
// Parameters:
//   - objectHashHex: SHA-256 hash of the object (64 lowercase hex)
//   - canonicalJSON: Optional canonical JSON bytes for audit replay
//
// Returns error if:
//   - objectHashHex is not valid 64-char lowercase hex
//   - File I/O fails
func AppendRegister(objectHashHex string, canonicalJSON []byte) error {
	// Validate object hash
	if !hex64Pattern.MatchString(objectHashHex) {
		return fmt.Errorf("%w: object_hash_hex must be 64 lowercase hex chars, got %q", ErrInvalidHex, objectHashHex)
	}

	// Create register entry
	entry := RegisterEntry{
		Type:          "register",
		Canon:         "v1.0",
		Timestamp:     time.Now().UTC().Format(time.RFC3339Nano),
		ObjectHashHex: objectHashHex,
	}

	// Optionally encode canonical JSON
	if len(canonicalJSON) > 0 {
		entry.CanonicalJSONB64 = base64.StdEncoding.EncodeToString(canonicalJSON)
	}

	return appendEntry(entry)
}

// ListRegistersSince returns all registration entries after the specified timestamp.
//
// Parameters:
//   - lastSealTS: Return only entries with timestamp > lastSealTS
//
// Returns:
//   - Slice of RegisterEntry records
//   - Error if ledger is corrupt or I/O fails
func ListRegistersSince(lastSealTS time.Time) ([]RegisterEntry, error) {
	path := GetLedgerPath()

	// If ledger doesn't exist, return empty slice
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []RegisterEntry{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to open ledger: %v", ErrLedgerIO, err)
	}
	defer file.Close()

	var registers []RegisterEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Parse as generic entry to determine type
		var entry struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("%w: line %d: invalid JSON: %v", ErrLedgerCorrupt, lineNum, err)
		}

		// Only process register entries
		if entry.Type == "register" {
			var reg RegisterEntry
			if err := json.Unmarshal(line, &reg); err != nil {
				return nil, fmt.Errorf("%w: line %d: invalid register entry: %v", ErrLedgerCorrupt, lineNum, err)
			}

			// Parse timestamp
			ts, err := time.Parse(time.RFC3339Nano, reg.Timestamp)
			if err != nil {
				return nil, fmt.Errorf("%w: line %d: invalid timestamp: %v", ErrLedgerCorrupt, lineNum, err)
			}

			// Filter by timestamp
			if ts.After(lastSealTS) {
				registers = append(registers, reg)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: failed to read ledger: %v", ErrLedgerIO, err)
	}

	return registers, nil
}

// AppendSeal appends a seal entry to the ledger.
//
// Parameters:
//   - manifest: Manifest containing Merkle root, signature, public key, and timestamp
//
// Returns error if:
//   - No registrations exist since last seal (or ever)
//   - Manifest validation fails
//   - File I/O fails
func AppendSeal(manifest Manifest) error {
	// Validate manifest fields
	if !hex64Pattern.MatchString(manifest.MerkleRoot) {
		return fmt.Errorf("%w: merkle_root must be 64 lowercase hex chars, got %q", ErrInvalidHex, manifest.MerkleRoot)
	}

	if !hex128Pattern.MatchString(manifest.Signature) {
		return fmt.Errorf("%w: signature must be 128 lowercase hex chars, got %q", ErrInvalidHex, manifest.Signature)
	}

	if !hex64Pattern.MatchString(manifest.PublicKey) {
		return fmt.Errorf("%w: public_key must be 64 lowercase hex chars, got %q", ErrInvalidHex, manifest.PublicKey)
	}

	// Validate timestamp format
	if _, err := time.Parse(time.RFC3339Nano, manifest.Timestamp); err != nil {
		return fmt.Errorf("%w: manifest timestamp: %v", ErrInvalidTimestamp, err)
	}

	// Check if there are any registrations to seal
	lastSealTS, err := getLastSealTimestamp()
	if err != nil {
		return err
	}

	registers, err := ListRegistersSince(lastSealTS)
	if err != nil {
		return err
	}

	if len(registers) == 0 {
		return ErrNoRegistrations
	}

	// Create seal entry
	entry := SealEntry{
		Type:     "seal",
		Manifest: manifest,
	}

	return appendEntry(entry)
}

// getLastSealTimestamp returns the timestamp of the last seal entry.
// Returns zero time if no seals exist.
func getLastSealTimestamp() (time.Time, error) {
	path := GetLedgerPath()

	// If ledger doesn't exist, return zero time
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return time.Time{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to open ledger: %v", ErrLedgerIO, err)
	}
	defer file.Close()

	var lastSealTS time.Time
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		var entry struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &entry); err != nil {
			return time.Time{}, fmt.Errorf("%w: line %d: invalid JSON: %v", ErrLedgerCorrupt, lineNum, err)
		}

		if entry.Type == "seal" {
			var seal SealEntry
			if err := json.Unmarshal(line, &seal); err != nil {
				return time.Time{}, fmt.Errorf("%w: line %d: invalid seal entry: %v", ErrLedgerCorrupt, lineNum, err)
			}

			ts, err := time.Parse(time.RFC3339Nano, seal.Manifest.Timestamp)
			if err != nil {
				return time.Time{}, fmt.Errorf("%w: line %d: invalid seal timestamp: %v", ErrLedgerCorrupt, lineNum, err)
			}

			lastSealTS = ts
		}
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("%w: failed to read ledger: %v", ErrLedgerIO, err)
	}

	return lastSealTS, nil
}

// appendEntry appends a JSON entry to the ledger file
func appendEntry(entry interface{}) error {
	ledgerMutex.Lock()
	defer ledgerMutex.Unlock()

	path := ledgerPath

	// Ensure ledger directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w: failed to create ledger directory: %v", ErrLedgerIO, err)
	}

	// Marshal entry to JSON
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("%w: failed to marshal entry: %v", ErrLedgerIO, err)
	}

	// Append to ledger file
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%w: failed to open ledger: %v", ErrLedgerIO, err)
	}
	defer file.Close()

	// Write JSON line
	if _, err := file.Write(jsonBytes); err != nil {
		return fmt.Errorf("%w: failed to write entry: %v", ErrLedgerIO, err)
	}

	// Write newline
	if _, err := file.Write([]byte("\n")); err != nil {
		return fmt.Errorf("%w: failed to write newline: %v", ErrLedgerIO, err)
	}

	return nil
}

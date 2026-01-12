package ledger

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// setupTestLedger creates a temporary ledger file for testing
func setupTestLedger(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	ledgerFile := filepath.Join(tempDir, "ledger.jsonl")
	SetLedgerPath(ledgerFile)
	return ledgerFile
}

// validObjectHash returns a valid 64-char lowercase hex string
func validObjectHash() string {
	return "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"
}

// validManifest returns a valid Manifest for testing
func validManifest() Manifest {
	return Manifest{
		MerkleRoot: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		Signature:  "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		PublicKey:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func TestAppendRegister_Valid(t *testing.T) {
	setupTestLedger(t)

	hash := validObjectHash()
	canonicalJSON := []byte(`{"key":"value"}`)

	err := AppendRegister(hash, canonicalJSON)
	if err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	// Verify entry was written
	data, err := os.ReadFile(GetLedgerPath())
	if err != nil {
		t.Fatalf("failed to read ledger: %v", err)
	}

	var entry RegisterEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if entry.Type != "register" {
		t.Errorf("Type = %s, want 'register'", entry.Type)
	}

	if entry.Canon != "v1.0" {
		t.Errorf("Canon = %s, want 'v1.0'", entry.Canon)
	}

	if entry.ObjectHashHex != hash {
		t.Errorf("ObjectHashHex = %s, want %s", entry.ObjectHashHex, hash)
	}

	// Verify canonical JSON was base64 encoded
	decoded, err := base64.StdEncoding.DecodeString(entry.CanonicalJSONB64)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	if string(decoded) != string(canonicalJSON) {
		t.Errorf("decoded JSON = %s, want %s", decoded, canonicalJSON)
	}

	// Verify timestamp is valid RFC3339Nano
	if _, err := time.Parse(time.RFC3339Nano, entry.Timestamp); err != nil {
		t.Errorf("invalid timestamp: %v", err)
	}
}

func TestAppendRegister_WithoutCanonicalJSON(t *testing.T) {
	setupTestLedger(t)

	hash := validObjectHash()

	err := AppendRegister(hash, nil)
	if err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	data, err := os.ReadFile(GetLedgerPath())
	if err != nil {
		t.Fatalf("failed to read ledger: %v", err)
	}

	var entry RegisterEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if entry.CanonicalJSONB64 != "" {
		t.Errorf("CanonicalJSONB64 should be empty, got %s", entry.CanonicalJSONB64)
	}
}

func TestAppendRegister_InvalidHex(t *testing.T) {
	setupTestLedger(t)

	tests := []struct {
		name string
		hash string
	}{
		{
			name: "too short",
			hash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae",
		},
		{
			name: "too long",
			hash: "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae300",
		},
		{
			name: "uppercase",
			hash: "A665A45920422F9D417E4867EFDC4FB8A04A1F3FFF1FA07E998E86F7F7A27AE3",
		},
		{
			name: "non-hex characters",
			hash: "g665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		},
		{
			name: "empty",
			hash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendRegister(tt.hash, nil)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "invalid hex format") {
				t.Errorf("expected ErrInvalidHex, got: %v", err)
			}
		})
	}
}

func TestListRegistersSince_EmptyLedger(t *testing.T) {
	setupTestLedger(t)

	registers, err := ListRegistersSince(time.Time{})
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(registers))
	}
}

func TestListRegistersSince_FilterByTimestamp(t *testing.T) {
	setupTestLedger(t)

	// Append three registers at different times
	hash1 := validObjectHash()
	if err := AppendRegister(hash1, nil); err != nil {
		t.Fatalf("AppendRegister 1 failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	cutoffTime := time.Now().UTC()
	time.Sleep(10 * time.Millisecond)

	hash2 := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if err := AppendRegister(hash2, nil); err != nil {
		t.Fatalf("AppendRegister 2 failed: %v", err)
	}

	hash3 := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if err := AppendRegister(hash3, nil); err != nil {
		t.Fatalf("AppendRegister 3 failed: %v", err)
	}

	// List registers since cutoff (should get 2 and 3, not 1)
	registers, err := ListRegistersSince(cutoffTime)
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 2 {
		t.Errorf("expected 2 registers, got %d", len(registers))
	}

	// Verify we got hash2 and hash3
	if registers[0].ObjectHashHex != hash2 {
		t.Errorf("registers[0] = %s, want %s", registers[0].ObjectHashHex, hash2)
	}
	if registers[1].ObjectHashHex != hash3 {
		t.Errorf("registers[1] = %s, want %s", registers[1].ObjectHashHex, hash3)
	}
}

func TestListRegistersSince_IgnoreSeals(t *testing.T) {
	setupTestLedger(t)

	// Append a register
	hash := validObjectHash()
	if err := AppendRegister(hash, nil); err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	// Append a seal
	manifest := validManifest()
	if err := AppendSeal(manifest); err != nil {
		t.Fatalf("AppendSeal failed: %v", err)
	}

	// List all registers (should only get the register, not the seal)
	registers, err := ListRegistersSince(time.Time{})
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 1 {
		t.Errorf("expected 1 register, got %d", len(registers))
	}

	if registers[0].ObjectHashHex != hash {
		t.Errorf("register hash = %s, want %s", registers[0].ObjectHashHex, hash)
	}
}

func TestListRegistersSince_CorruptLine(t *testing.T) {
	ledgerPath := setupTestLedger(t)

	// Write a valid entry
	hash := validObjectHash()
	if err := AppendRegister(hash, nil); err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	// Manually append corrupt line
	file, err := os.OpenFile(ledgerPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open ledger: %v", err)
	}
	file.WriteString("this is not valid json\n")
	file.Close()

	// ListRegistersSince should fail
	_, err = ListRegistersSince(time.Time{})
	if err == nil {
		t.Fatalf("expected error for corrupt ledger, got nil")
	}

	if !strings.Contains(err.Error(), "ledger corrupt") {
		t.Errorf("expected ErrLedgerCorrupt, got: %v", err)
	}
}

func TestAppendSeal_NoRegistrations(t *testing.T) {
	setupTestLedger(t)

	manifest := validManifest()

	err := AppendSeal(manifest)
	if err == nil {
		t.Fatalf("expected error when no registrations, got nil")
	}

	if !strings.Contains(err.Error(), "no registrations to seal") {
		t.Errorf("expected ErrNoRegistrations, got: %v", err)
	}
}

func TestAppendSeal_Valid(t *testing.T) {
	setupTestLedger(t)

	// First append a register
	hash := validObjectHash()
	if err := AppendRegister(hash, nil); err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	// Now seal
	manifest := validManifest()
	if err := AppendSeal(manifest); err != nil {
		t.Fatalf("AppendSeal failed: %v", err)
	}

	// Verify seal was written
	data, err := os.ReadFile(GetLedgerPath())
	if err != nil {
		t.Fatalf("failed to read ledger: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// Parse seal entry (second line)
	var seal SealEntry
	if err := json.Unmarshal([]byte(lines[1]), &seal); err != nil {
		t.Fatalf("failed to unmarshal seal: %v", err)
	}

	if seal.Type != "seal" {
		t.Errorf("Type = %s, want 'seal'", seal.Type)
	}

	if seal.Manifest.MerkleRoot != manifest.MerkleRoot {
		t.Errorf("MerkleRoot = %s, want %s", seal.Manifest.MerkleRoot, manifest.MerkleRoot)
	}

	if seal.Manifest.Signature != manifest.Signature {
		t.Errorf("Signature = %s, want %s", seal.Manifest.Signature, manifest.Signature)
	}

	if seal.Manifest.PublicKey != manifest.PublicKey {
		t.Errorf("PublicKey = %s, want %s", seal.Manifest.PublicKey, manifest.PublicKey)
	}
}

func TestAppendSeal_InvalidManifest(t *testing.T) {
	setupTestLedger(t)

	// Append a register first
	if err := AppendRegister(validObjectHash(), nil); err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	tests := []struct {
		name     string
		manifest Manifest
		errMsg   string
	}{
		{
			name: "invalid merkle root",
			manifest: Manifest{
				MerkleRoot: "invalid",
				Signature:  "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				PublicKey:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			},
			errMsg: "invalid hex format",
		},
		{
			name: "invalid signature length",
			manifest: Manifest{
				MerkleRoot: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				Signature:  "short",
				PublicKey:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			},
			errMsg: "invalid hex format",
		},
		{
			name: "invalid public key",
			manifest: Manifest{
				MerkleRoot: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				Signature:  "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				PublicKey:  "INVALID",
				Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			},
			errMsg: "invalid hex format",
		},
		{
			name: "invalid timestamp",
			manifest: Manifest{
				MerkleRoot: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
				Signature:  "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				PublicKey:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				Timestamp:  "not-a-timestamp",
			},
			errMsg: "invalid timestamp format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AppendSeal(tt.manifest)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing %q, got: %v", tt.errMsg, err)
			}
		})
	}
}

func TestMultipleEntries(t *testing.T) {
	setupTestLedger(t)

	// Append multiple registers
	for i := 0; i < 5; i++ {
		hash := validObjectHash()
		if err := AppendRegister(hash, nil); err != nil {
			t.Fatalf("AppendRegister %d failed: %v", i, err)
		}
	}

	// Seal
	if err := AppendSeal(validManifest()); err != nil {
		t.Fatalf("AppendSeal failed: %v", err)
	}

	// Append more registers
	for i := 0; i < 3; i++ {
		hash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
		if err := AppendRegister(hash, nil); err != nil {
			t.Fatalf("AppendRegister after seal %d failed: %v", i, err)
		}
	}

	// List all registers
	registers, err := ListRegistersSince(time.Time{})
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 8 {
		t.Errorf("expected 8 registers, got %d", len(registers))
	}
}

func TestAppendSeal_AfterPreviousSeal(t *testing.T) {
	setupTestLedger(t)

	// First round: register + seal
	if err := AppendRegister(validObjectHash(), nil); err != nil {
		t.Fatalf("AppendRegister 1 failed: %v", err)
	}
	if err := AppendSeal(validManifest()); err != nil {
		t.Fatalf("AppendSeal 1 failed: %v", err)
	}

	// Try to seal again without new registrations
	err := AppendSeal(validManifest())
	if err == nil {
		t.Fatalf("expected error when sealing without new registrations")
	}
	if !strings.Contains(err.Error(), "no registrations to seal") {
		t.Errorf("expected ErrNoRegistrations, got: %v", err)
	}

	// Add new registration
	hash2 := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if err := AppendRegister(hash2, nil); err != nil {
		t.Fatalf("AppendRegister 2 failed: %v", err)
	}

	// Now seal should work
	if err := AppendSeal(validManifest()); err != nil {
		t.Fatalf("AppendSeal 2 failed: %v", err)
	}
}

func TestCanonicalJSONB64_RoundTrip(t *testing.T) {
	setupTestLedger(t)

	originalJSON := []byte(`{"user":"alice","score":100,"active":true}`)
	hash := validObjectHash()

	if err := AppendRegister(hash, originalJSON); err != nil {
		t.Fatalf("AppendRegister failed: %v", err)
	}

	registers, err := ListRegistersSince(time.Time{})
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 1 {
		t.Fatalf("expected 1 register, got %d", len(registers))
	}

	decoded, err := base64.StdEncoding.DecodeString(registers[0].CanonicalJSONB64)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	if string(decoded) != string(originalJSON) {
		t.Errorf("round-trip failed: got %s, want %s", decoded, originalJSON)
	}
}

func TestConcurrentAppends(t *testing.T) {
	setupTestLedger(t)

	// Test that concurrent appends don't corrupt the ledger
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 5; j++ {
				hash := validObjectHash()
				if err := AppendRegister(hash, nil); err != nil {
					t.Errorf("concurrent AppendRegister failed: %v", err)
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify ledger integrity
	registers, err := ListRegistersSince(time.Time{})
	if err != nil {
		t.Fatalf("ListRegistersSince failed: %v", err)
	}

	if len(registers) != 50 {
		t.Errorf("expected 50 registers, got %d", len(registers))
	}
}

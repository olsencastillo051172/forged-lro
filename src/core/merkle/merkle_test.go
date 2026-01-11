package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
)

// leafHash: helper para producir hojas válidas (SHA-256 hex lowercase) desde strings.
func leafHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func makeLeaves(vals []string) []string {
	out := make([]string, len(vals))
	for i, v := range vals {
		out[i] = leafHash(v)
	}
	return out
}

// -------------------------
// Sanidad (conservar)
// -------------------------

func TestVerifyProof_EmptyProofWhenTotalLeavesGreaterThanOne(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	ok, err := VerifyProof(leaves[0], 0, len(leaves), []ProofNode{}, root)
	if err == nil || !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof, got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}

func TestVerifyProof_InvalidProofHash(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	proof, _, err := BuildProof(leaves, 0)
	if err != nil {
		t.Fatalf("BuildProof error: %v", err)
	}
	if len(proof) == 0 {
		t.Fatalf("expected non-empty proof")
	}

	// Corromper hash del primer nodo a un formato inválido (no 64-hex)
	proof[0].Hash = "zzz"

	ok, err := VerifyProof(leaves[0], 0, len(leaves), proof, root)
	if err == nil || !errors.Is(err, ErrInvalidLeafFormat) {
		t.Fatalf("expected ErrInvalidLeafFormat, got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}

func TestVerifyProof_SingleLeafWithNonEmptyProof(t *testing.T) {
	leaves := makeLeaves([]string{"A"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	// Proof no vacío con totalLeaves==1 => MALFORMED
	proof := []ProofNode{{Hash: leaves[0], Position: "left"}}

	ok, err := VerifyProof(leaves[0], 0, len(leaves), proof, root)
	if err == nil || !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof, got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}

// -------------------------
// Huecos críticos (cerrar)
// -------------------------

func TestVerifyProof_InvalidPositionRejected(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B", "C", "D"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	proof, _, err := BuildProof(leaves, 2) // hoja "C"
	if err != nil {
		t.Fatalf("BuildProof error: %v", err)
	}
	if len(proof) == 0 {
		t.Fatalf("expected non-empty proof")
	}

	// Corromper posición del primer paso
	if proof[0].Position == "left" {
		proof[0].Position = "right"
	} else {
		proof[0].Position = "left"
	}

	ok, err := VerifyProof(leaves[2], 2, len(leaves), proof, root)
	if err == nil || !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof for invalid position; got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}

func TestVerifyProof_WrongLengthRejected(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B", "C", "D"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	proof, _, err := BuildProof(leaves, 1)
	if err != nil {
		t.Fatalf("BuildProof error: %v", err)
	}
	if len(proof) < 2 {
		t.Fatalf("expected proof length >= 2 for 4 leaves, got %d", len(proof))
	}

	// Recortar proof (length inválida)
	proof = proof[:len(proof)-1]

	ok, err := VerifyProof(leaves[1], 1, len(leaves), proof, root)
	if err == nil || !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof for wrong length; got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}

func TestVerifyProof_OddDuplicationRuleBound(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B", "C"}) // 3 hojas (impar)
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	// Último índice para forzar duplicación
	idx := 2
	proof, _, err := BuildProof(leaves, idx)
	if err != nil {
		t.Fatalf("BuildProof error: %v", err)
	}
	if len(proof) == 0 {
		t.Fatalf("expected non-empty proof")
	}

	// En el primer nivel (n=3), idx=2 (par) => sibling duplicado => proof[0].Hash debe == leaf
	// Corromperlo a otro hash válido
	proof[0].Hash = leaves[1]

	ok, err := VerifyProof(leaves[idx], idx, len(leaves), proof, root)
	if err == nil || !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof for odd-duplication violation; got ok=%v err=%v", ok, err)
	}
	if ok {
		t.Fatalf("expected ok=false; got ok=true")
	}
}



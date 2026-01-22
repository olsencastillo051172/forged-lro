package policy

import (
	"bytes"
	"testing"
)

func TestCanonicalizePolicy_Determinism(t *testing.T) {
	// Definimos dos políticas semánticamente idénticas
	p1 := &RotationPolicy{
		PolicyVersion: "1.0",
		Issuer: IssuerInfo{Name: "Alpha", ID: "rva://1"},
	}
	
	p2 := &RotationPolicy{
		// Invertimos el orden lógico de asignación (aunque en structs no importa, 
		// el test asegura que el marshaller sea determinista)
		Issuer: IssuerInfo{ID: "rva://1", Name: "Alpha"},
		PolicyVersion: "1.0",
	}

	b1, err := CanonicalizePolicy(p1)
	if err != nil {
		t.Fatalf("Failed to canonicalize p1: %v", err)
	}

	b2, err := CanonicalizePolicy(p2)
	if err != nil {
		t.Fatalf("Failed to canonicalize p2: %v", err)
	}

	if !bytes.Equal(b1, b2) {
		t.Errorf("Determinism failed: bytes are not identical\nB1: %s\nB2: %s", b1, b2)
	}

	if bytes.HasSuffix(b1, []byte("\n")) {
		t.Error("Format error: canonical bytes must not have a trailing newline")
	}
}

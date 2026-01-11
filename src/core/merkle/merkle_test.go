func TestVerifyProof_EmptyProofWhenTotalLeavesGreaterThanOne(t *testing.T) {
	leaves := makeLeaves([]string{"A", "B"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	// Proof vacío con totalLeaves=2 => debe ser MALFORMED => err != nil y ok == false
	ok, err := VerifyProof(leaves[0], 0, len(leaves), []ProofNode{}, root)
	if err == nil {
		t.Fatalf("expected error for empty proof when totalLeaves > 1; got ok=%v err=nil", ok)
	}
	if ok {
		t.Fatalf("expected ok=false for empty proof when totalLeaves > 1; got ok=true err=%v", err)
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

	// Corromper hash del primer nodo => debe ser MALFORMED => err != nil y ok == false
	if len(proof) == 0 {
		t.Fatalf("expected non-empty proof for totalLeaves=2")
	}

	proof[0].Hash = "zzz"
	ok, err := VerifyProof(leaves[0], 0, len(leaves), proof, root)
	if err == nil {
		t.Fatalf("expected error for invalid proof node hash; got ok=%v err=nil", ok)
	}
	if ok {
		t.Fatalf("expected ok=false for invalid proof node hash; got ok=true err=%v", err)
	}
}

func TestVerifyProof_SingleLeafWithNonEmptyProof(t *testing.T) {
	leaves := makeLeaves([]string{"A"})
	root, err := BuildRoot(leaves)
	if err != nil {
		t.Fatalf("BuildRoot error: %v", err)
	}

	// Proof no vacío aunque totalLeaves=1 => debe ser MALFORMED => err != nil y ok == false
	proof := []ProofNode{{Hash: leaves[0], Position: "left"}}
	ok, err := VerifyProof(leaves[0], 0, len(leaves), proof, root)
	if err == nil {
		t.Fatalf("expected error for non-empty proof when totalLeaves == 1; got ok=%v err=nil", ok)
	}
	if ok {
		t.Fatalf("expected ok=false for non-empty proof when totalLeaves == 1; got ok=true err=%v", err)
	}
}


package merkle

import (
    "testing"
)

// Benchmark construcción de árbol con 1,000 hojas
func BenchmarkBuildTree1000(b *testing.B) {
    leaves := make([]HashHex, 1000)
    for i := 0; i < 1000; i++ {
        leaves[i] = Sha256String(string(rune(i)))
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := Build(leaves)
        if err != nil {
            b.Fatalf("Build failed: %v", err)
        }
    }
}

// Benchmark verificación de prueba con profundidad ~20 (árbol ~1 millón de nodos)
func BenchmarkVerifyDepth20(b *testing.B) {
    // Generamos ~1 millón de hojas
    leaves := make([]HashHex, 1<<20) // 2^20 ≈ 1,048,576
    for i := 0; i < len(leaves); i++ {
        leaves[i] = Sha256String(string(rune(i)))
    }
    tree, err := Build(leaves)
    if err != nil {
        b.Fatalf("Build failed: %v", err)
    }

    // Seleccionamos una hoja en el medio
    idx := len(leaves) / 2
    proof, err := GenerateProof(leaves, idx)
    if err != nil {
        b.Fatalf("Proof generation failed: %v", err)
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        res, err := Verify(proof, tree.Root)
        if err != nil {
            b.Fatalf("Verify failed: %v", err)
        }
        if !res.IsValid {
            b.Fatalf("Proof invalid")
        }
    }
}

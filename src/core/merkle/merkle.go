package merkle

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "regexp"
)

var (
    // ErrEmptyLeaves is returned when attempting to build a tree with no leaves.
    ErrEmptyLeaves = errors.New("cannot build merkle tree from empty leaf set")

    // ErrInvalidLeafFormat is returned when a leaf hash doesn't match the required format.
    ErrInvalidLeafFormat = errors.New("leaf hash must be 64-character lowercase hex string")

    // ErrInvalidIndex is returned when proof index is out of bounds.
    ErrInvalidIndex = errors.New("index out of bounds")

    // ErrInvalidProof is returned when proof verification fails due to invalid structure.
    ErrInvalidProof = errors.New("invalid merkle proof")

    // ErrInvalidTotalLeaves is returned when totalLeaves is invalid.
    ErrInvalidTotalLeaves = errors.New("invalid totalLeaves parameter")
)

// hashPattern validates SHA-256 hex strings (64 lowercase hex chars).
var hashPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

// ProofNode represents a single node in a Merkle proof path.
// Position indicates whether this hash should be concatenated on the "left" or "right"
// when reconstructing the path to the root.
type ProofNode struct {
    Hash     string `json:"hash"`
    Position string `json:"position"` // "left" or "right"
}

// BuildRoot constructs the Merkle root from an ordered list of leaf hashes.
//
// Canon rules:
//   - Leaves are SHA-256 hashes (64-char lowercase hex), validated before processing
//   - Tree is built in order; leaves are NOT sorted
//   - Parent hash = SHA256(left_bytes || right_bytes) where || is byte concatenation
//   - If odd number of nodes at any level, duplicate the last node to form a pair
//   - Single leaf returns that leaf as root (no extra hashing)
//   - Empty leaf set returns error
func BuildRoot(leaves []string) (string, error) {
    if len(leaves) == 0 {
        return "", ErrEmptyLeaves
    }

    for i, leaf := range leaves {
        if !hashPattern.MatchString(leaf) {
            return "", fmt.Errorf("%w: leaf[%d] = %q", ErrInvalidLeafFormat, i, leaf)
        }
    }

    // Single leaf: root is the leaf itself (no hashing).
    if len(leaves) == 1 {
        return leaves[0], nil
    }

    currentLevel := make([]string, len(leaves))
    copy(currentLevel, leaves)

    for len(currentLevel) > 1 {
        nextLevel := make([]string, 0, (len(currentLevel)+1)/2)

        for i := 0; i < len(currentLevel); i += 2 {
            left := currentLevel[i]
            right := left
            if i+1 < len(currentLevel) {
                right = currentLevel[i+1]
            }

            parent, err := hashPair(left, right)
            if err != nil {
                return "", err
            }
            nextLevel = append(nextLevel, parent)
        }

        currentLevel = nextLevel
    }

    return currentLevel[0], nil
}

// BuildProof generates a Merkle proof for the leaf at the specified index.
//
// The proof consists of sibling hashes along the path from leaf to root.
// Each ProofNode indicates whether the sibling should be on the "left" or "right"
// when reconstructing parent hashes.
func BuildProof(leaves []string, index int) ([]ProofNode, string, error) {
    if len(leaves) == 0 {
        return nil, "", ErrEmptyLeaves
    }

    if index < 0 || index >= len(leaves) {
        return nil, "", fmt.Errorf("%w: index %d, total leaves %d", ErrInvalidIndex, index, len(leaves))
    }

    for i, leaf := range leaves {
        if !hashPattern.MatchString(leaf) {
            return nil, "", fmt.Errorf("%w: leaf[%d] = %q", ErrInvalidLeafFormat, i, leaf)
        }
    }

    // Single leaf: no proof needed, root is the leaf.
    if len(leaves) == 1 {
        return []ProofNode{}, leaves[0], nil
    }

    proof := make([]ProofNode, 0, 32)
    currentLevel := make([]string, len(leaves))
    copy(currentLevel, leaves)
    currentIndex := index

    for len(currentLevel) > 1 {
        nextLevel := make([]string, 0, (len(currentLevel)+1)/2)

        for i := 0; i < len(currentLevel); i += 2 {
            left := currentLevel[i]
            right := left
            if i+1 < len(currentLevel) {
                right = currentLevel[i+1]
            }

            // If this pair contains our target index, record sibling.
            if i == currentIndex || i+1 == currentIndex {
                if currentIndex%2 == 0 {
                    // Target is left child, sibling is right.
                    proof = append(proof, ProofNode{Hash: right, Position: "right"})
                } else {
                    // Target is right child, sibling is left.
                    proof = append(proof, ProofNode{Hash: left, Position: "left"})
                }
            }

            parent, err := hashPair(left, right)
            if err != nil {
                return nil, "", err
            }
            nextLevel = append(nextLevel, parent)
        }

        currentIndex = currentIndex / 2
        currentLevel = nextLevel
    }

    return proof, currentLevel[0], nil
}

// VerifyProof verifies a Merkle proof without rebuilding the entire tree.
//
// Returns true if proof is valid, false otherwise.
// Returns error for invalid inputs or malformed proof structure.
func VerifyProof(leaf string, index int, totalLeaves int, proof []ProofNode, expectedRoot string) (bool, error) {
    if !hashPattern.MatchString(leaf) {
        return false, fmt.Errorf("%w: leaf = %q", ErrInvalidLeafFormat, leaf)
    }
    if !hashPattern.MatchString(expectedRoot) {
        return false, fmt.Errorf("%w: expectedRoot = %q", ErrInvalidLeafFormat, expectedRoot)
    }

    if totalLeaves <= 0 {
        return false, fmt.Errorf("%w: totalLeaves must be positive", ErrInvalidTotalLeaves)
    }
    if index < 0 || index >= totalLeaves {
        return false, fmt.Errorf("%w: index %d, totalLeaves %d", ErrInvalidIndex, index, totalLeaves)
    }

    // Single leaf: proof must be empty and leaf must equal root.
    if totalLeaves == 1 {
        if len(proof) != 0 {
            return false, fmt.Errorf("%w: single leaf should have empty proof", ErrInvalidProof)
        }
        return leaf == expectedRoot, nil
    }

    // For totalLeaves > 1, proof should not be empty (minimal sanity).
    if len(proof) == 0 {
        return false, fmt.Errorf("%w: proof is empty but totalLeaves > 1", ErrInvalidProof)
    }

    for i, node := range proof {
        if !hashPattern.MatchString(node.Hash) {
            return false, fmt.Errorf("%w: proof[%d].hash = %q", ErrInvalidLeafFormat, i, node.Hash)
        }
        if node.Position != "left" && node.Position != "right" {
            return false, fmt.Errorf("%w: proof[%d].position must be 'left' or 'right', got %q", ErrInvalidProof, i, node.Position)
        }
    }

    // Reconstruct root from leaf using proof path.
    currentHash := leaf
    for _, node := range proof {
        var left, right string
        if node.Position == "right" {
            left = currentHash
            right = node.Hash
        } else {
            left = node.Hash
            right = currentHash
        }

        parent, err := hashPair(left, right)
        if err != nil {
            return false, err
        }
        currentHash = parent
    }

    return currentHash == expectedRoot, nil
}

// hashPair combines two hex-encoded hashes into a parent hash.
// Decodes both hashes to bytes, concatenates left||right, applies SHA-256, returns lowercase hex.
func hashPair(leftHex, rightHex string) (string, error) {
    leftBytes, err := hex.DecodeString(leftHex)
    if err != nil {
        return "", fmt.Errorf("failed to decode left hash: %w", err)
    }
    rightBytes, err := hex.DecodeString(rightHex)
    if err != nil {
        return "", fmt.Errorf("failed to decode right hash: %w", err)
    }

    combined := append(leftBytes, rightBytes...) // left||right (bytes)
    sum := sha256.Sum256(combined)
    return hex.EncodeToString(sum[:]), nil
}

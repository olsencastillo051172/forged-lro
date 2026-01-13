package sign

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// DeriveKeyPairFromSeedHex derives Ed25519 keypair deterministically from a 32-byte seed hex (64 hex chars).
// Returns:
// - privHex: 64-byte private key hex (128 chars) from ed25519.NewKeyFromSeed
// - pubHex:  32-byte public key hex (64 chars)
func DeriveKeyPairFromSeedHex(seedHex string) (privHex string, pubHex string, err error) {
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return "", "", fmt.Errorf("invalid seed hex: %w", err)
	}
	if len(seed) != 32 {
		return "", "", fmt.Errorf("seed must be 32 bytes")
	}

	priv := ed25519.NewKeyFromSeed(seed) // 64 bytes
	pub := priv.Public().(ed25519.PublicKey)

	return hex.EncodeToString(priv), hex.EncodeToString(pub), nil
}

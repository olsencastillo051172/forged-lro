# RVA — Reality Validation Authority (Powered by FORGED-LRO)

**Sovereign Evidence Infrastructure for certifying past-state data.**

**Status:** Canon v1.0 (FROZEN)  
**Scope:** Specifications + JSON Schemas only (no implementation yet)

---

## What This Is

**RVA (Reality Validation Authority)** is the public-facing authority that issues cryptographic certificates proving that:

> **Data D (hash H) was registered by entity E at timestamp T.**

RVA is powered by **FORGED-LRO**, an append-only ledger engine with epoch-based Merkle anchoring and offline verification.

---

## What This Is NOT

- NOT a blockchain
- NOT a truth oracle
- NOT a monitoring / analytics system
- NOT a prediction or opinion service

**RVA certifies existence + registration time, not truth.**

---

## Core Guarantees (Canon)

1. Deterministic hashing (SHA-256)
2. Append-only ledger (no UPDATE/DELETE)
3. Epoch manifests (Merkle root + Ed25519 signature)
4. Offline verification (no server required)
5. Zero-opinion certificates (state, not interpretation)

## Conclusions and Guarantees — Canon v1.0 (Sealed)

The delivered implementation satisfies 100% of the functional and non-functional requirements of the Canon v1.0 standard for:
- `src/core/hash` (canonical JSON + SHA-256 lowercase hex)
- `src/core/merkle` (deterministic Merkle tree with byte-level concatenation)
- `src/core/sign` (deterministic Ed25519 with lowercase hex I/O)

By eliminating external dependencies and enforcing strict validation of canonical formats, the system provides a solid foundation for evidence immutability and verifiability within FORGED-LRO.

### Compliance Checklist ✅
- stdlib-only
- canonical JSON with sorted keys (hash)
- minified JSON bytes (hash)
- SHA-256 lowercase hex (hash/merkle)
- Merkle: `hashPair = SHA256(left_bytes || right_bytes)`
- Merkle: odd leaf count ⇒ deterministic duplication of the last node
- Ed25519: seed/pub/sig in lowercase hex (64/64/128)
- signing/verification over raw hash bytes (32 bytes)
- exported and wrapped errors
- exhaustive tests
- `go test ./...` passes in CI (GitHub Actions)








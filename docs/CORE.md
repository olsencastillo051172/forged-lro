# FORGED-LRO CORE — Canon v1.0

**Status:** FROZEN  
**Scope:** Core ledger engine specification (no implementation)  
**Role:** Deterministic engine powering RVA and all future verticals

---

## 1. Ledger Structure (Append-Only)

FORGED-LRO operates as an **append-only ledger**.

### Core Rules
- NO update
- NO delete
- NO rewrite
- NO retroactive mutation

Each entry is immutable and cryptographically chained.

### Logical Entry Fields
- `entry_id` (UUID)
- `vertical` (RVA_CORE | INVOICE | SUPPLY_CHAIN | future)
- `payload_hash` (SHA-256)
- `prev_entry_hash` (SHA-256 | null for genesis)
- `entry_hash` (SHA-256 of full entry)
- `submitter_id`
- `server_received_utc` (authoritative timestamp)
- `epoch_id`
- `entry_no`
- `metadata` (opaque to core)

Storage backend is intentionally undefined.

---

## 2. Epoch System

### 2.1 Definition

An **epoch** is a fixed-size group of ledger entries.

**Canon constant:**
- `EPOCH_SIZE = 1000`

### 2.2 Epoch Closure

When an epoch reaches size:
1. All `entry_hash` values are used as Merkle tree leaves
2. A deterministic Merkle tree is constructed
3. A single **Merkle root** is generated
4. The epoch is sealed and immutable

### 2.3 Epoch Manifest

Each epoch produces a manifest containing:
- `epoch_id`
- `entry_range`
- `merkle_root`
- `generation_timestamp`
- `rva_signature` (Ed25519)

Manifests are WORM artifacts.

---

## 3. Certificates (RVA Format)

FORGED-LRO does not issue certificates directly.  
It provides the cryptographic guarantees used by RVA.

### Canonical Statement

All certificates reduce to:

> **Data D (hash H) was registered by entity E at timestamp T**

### Mandatory Certificate Elements
- `payload_hash (H)`
- `submitter_id (E)`
- `registration_timestamp (T)`
- `ledger_entry_id`
- `epoch_id`
- `merkle_proof`
- `merkle_root`
- `rva_signature`

Exact schemas are defined in `/schemas`.

---

## 4. Verification Flow (Offline)

Verification requires **no API, no server, no trust**.

### Offline Steps
1. Recompute payload hash (SHA-256)
2. Validate Merkle path → Merkle root
3. Match Merkle root to epoch manifest
4. Verify RVA signature (Ed25519)
5. Confirm timestamps and identifiers

If all checks pass, the certificate is valid.

---

## 5. Interfaces (Minimal)

### 5.1 Submit Interface
Used by RVA and authorized submitters.

Responsibilities:
- Receive payload
- Hash deterministically
- Append ledger entry
- Return certificate material

### 5.2 Verification Interface
Used by anyone.

Capabilities:
- Verify Merkle proof
- Verify signature
- No authentication required

CLI/API specifics are implementation concerns.

---

## 6. Cryptographic Constants (Inherited from Canon)

- Hash: **SHA-256**
- Signature: **Ed25519**
- Time standard: **UTC only**
- Encoding: UTF-8
- Serialization: Canonical JSON

Deviation requires Canon v1.1.

---

## 7. Infrastructure Neutrality

FORGED-LRO is:
- Cloud-agnostic
- Multi-cloud compatible
- On-prem compatible
- Vendor-neutral

Infrastructure choices do not affect protocol validity.

---

## 8. Change Control

This document is frozen under Canon v1.0.

Any change requires:
- New Canon version
- Explicit migration rules
- No retroactive edits

---

FORGED-LRO CORE  
Deterministic. Append-only. Auditable. Sovereign.

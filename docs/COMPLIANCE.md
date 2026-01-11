# COMPLIANCE.md — Canon v1.0

## Purpose
This document defines the official compliance checklist for FORGED-LRO.  
It ensures that all objects conform to the six canonical schemas and that the system remains audit-ready.

---

## Verification Steps

### 1. Invoice Payload
- Validate against `/schemas/invoice_payload.schema.json`
- Required fields:
  - `document_hash` (SHA-256, 64 hex characters)
  - `submitter_id` (string)
  - `submission_timestamp` (UTC date-time)

### 2. Supply Chain Event
- Validate against `/schemas/supply_chain_event.schema.json`
- Required fields:
  - `event_id` (UUID)
  - `device_time_utc` (UTC date-time)
  - `device_id` (string)

### 3. Certificate
- Validate against `/schemas/certificate.schema.json`
- Required fields:
  - `certificate_type` (string)
  - `entry_id` (UUID)
  - `epoch_id` (integer)
  - `merkle_root` (SHA-256, 64 hex characters)
- Optional:
  - `rva_signature` (string)

### 4. Epoch Manifest
- Validate against `/schemas/epoch_manifest.schema.json`
- Required fields:
  - `epoch_id` (integer)
  - `merkle_root` (SHA-256, 64 hex characters)
  - `signature` (string)

### 5. Merkle Proof
- Validate against `/schemas/merkle_tree.schema.json`
- Each array item must contain:
  - `hash` (SHA-256, 64 hex characters)
  - `position` (enum: `left` or `right`)

### 6. Ledger Entry
- Validate against `/schemas/ledger_entry.schema.json`
- Required fields:
  - `entry_id` (UUID)
  - `payload_hash` (SHA-256, 64 hex characters)
  - `submitter_id` (string)
  - `timestamp_utc` (UTC date-time)
  - `epoch_id` (integer)
- Optional:
  - `metadata` (object)

---

## Compliance Result
- If all objects pass validation against their schemas → **Audit OK**  
- If discrepancies are found → **Record finding and correct before freezing Canon**

---

## Notes
- This file is a **single explanatory document**.  
- It is frozen together with `CANON.md` and `CORE.md` as part of **Canon v1.0**.  
- Any schema updates must be reflected here to maintain audit consistency.


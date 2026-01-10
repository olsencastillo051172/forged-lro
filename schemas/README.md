# RVA / FORGED-LRO Schemas — Canon v1.0

**Status:** FROZEN (Canon v1.0)  
**Purpose:** JSON Schema validation for all RVA inputs/outputs

---

## Files

1) `ledger_entry.schema.json`  
Universal append-only ledger envelope written by FORGED-LRO. Vertical payloads are hashed and referenced by this envelope.

2) `invoice_payload.schema.json`  
Invoice/document registration payload (supports exact file hashing and canonical fingerprinting).

3) `supply_chain_event.schema.json`  
Supply chain event payload with dual timestamps and evidence-preserving flags.

4) `certificate.schema.json`  
Universal certificate format issued by RVA, including Merkle proof and Ed25519 signature.

5) `epoch_manifest.schema.json`  
Epoch manifest emitted every `EPOCH_SIZE` entries, containing Merkle root + signing metadata.

6) `merkle_tree.schema.json`  
Optional full Merkle tree export rules for auditors.

---

## Change Control

Canon v1.0 schemas are frozen.  
Any modification requires Canon v1.1 + migration notes + approval.

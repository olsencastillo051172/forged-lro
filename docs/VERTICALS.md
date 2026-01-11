# RVA Verticals — Canon v1.0
**Status:** Frozen  

RVA supports multiple verticals, all powered by the same sovereign evidence core.

---

## 1. Vertical 1 — Document & PDF Validation
(Invoice Double-Spend / General Document Notarization)

### Purpose
Certify that a document existed at a specific time and was registered once.

### Use Cases
- Invoices (anti double-spend)
- Contracts
- Legal documents
- Financial statements

### Guarantee
"Document hash H was registered at timestamp T by entity E."

**Schema Reference:** `/schemas/document_certificate.json`

---

## 2. Vertical 2 — Supply Chain Event Certification

### Purpose
Certify physical-world events reported digitally.

### Use Cases
- Cold chain monitoring
- Custody transfers
- GPS checkpoints
- Tamper alerts

### Guarantee
"Event E occurred at location L and time T as reported by device D."

**Schema Reference:** `/schemas/supply_event.json`

---

## 3. Vertical 3 — Reality Certification (Core)

### Purpose
Certify existence of any digital artifact.

### Use Cases
- Datasets
- Media files
- API responses
- System states

### Guarantee
"Data D existed in state S at time T."

**Schema Reference:** `/schemas/reality_certificate.json`

---

## 4. Vertical Architecture Principle

- One ledger  
- One canonical format  
- Multiple vertical payloads  
- Uniform verification logic  

All verticals share:
- Same hashing rules  
- Same epoch logic  
- Same verification flow  

---

## 5. Expansion Policy

New verticals:
- MUST reuse the same ledger  
- MUST not alter existing guarantees  
- MUST define payload schemas only  

No vertical may redefine truth semantics.

---

## 6. Integrity Statement (Verticals)

Any vertical that:
- Alters canonical invariants, or  
- Introduces truth semantics, or  
- Deviates from the core guarantees  

**MUST NOT** claim RVA compliance.  
Compliance is binary.

---

## End of RVA Verticals Specification
**Canon v1.0 — Frozen**  
**Effective:** 2026-01-10

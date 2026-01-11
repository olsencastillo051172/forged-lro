# COMPLIANCE.md — FORGED-LRO Phase A

## Scope
- Canon v1.0 invariants frozen in `src/config/canon.go`
- Guardrail tests in `tests/canon_test.go`
- CLI scaffold in `cli/verify_certificate.go`
- Schemas in `/schemas/`
- No business logic included
- No retro-edits to `/docs` or `/schemas`

## Guardrails
- Constants are immutable (`CanonStatus = FROZEN`, `CanonVersion = v1.0`)
- Tests enforce invariants (`go test ./tests/...`)
- CLI interface defined but not implemented
- `.gitignore` excludes secrets, evidence packs, and local artifacts

## Audit Procedure
1. Clone repository
2. Run `go mod tidy`
3. Run `go test ./tests/...`
4. Verify all tests pass
5. Confirm commit history shows forward-only discipline (no retro-edits)

## Forward-Only Discipline
- Canon v1.0 is sealed and cannot be modified
- All changes are forward-only via new commits/PRs
- Audit trail preserved in Git history
- No retroactive edits allowed to invariants, schemas, or compliance docs

#!/bin/bash
# RVA Audit Script
# One-command pipeline: generate anchor (Go), run tests, inject into portal, build static dossier.
# Stops on error, outputs final artifact in rva-demo-portal/out.

set -e

LOGFILE="audit.log"
exec > >(tee -a "$LOGFILE") 2>&1

echo "⚔️ RVA Sovereign Audit: starting notarization at $(date -u +"%Y-%m-%dT%H:%M:%SZ")"

# 1. Generate Anchor with Go
echo "➡️ Generating anchor with Go..."
go run cmd/rva-rotate/main.go

# 2. Run Go tests
echo "🧪 Running Go tests..."
go test ./lib/rva/... -v

# 3. Inject Anchor into Portal
echo "➡️ Injecting anchor into portal..."
cd rva-demo-portal
npm ci
npm run inject-anchor

# 4. Run Portal tests
echo "🧪 Running Portal tests..."
npm test

# 5. Build Portal
echo "➡️ Building forensic portal..."
npm run build

# 6. Final message
echo "✅ Sovereign notarization complete."
echo "📂 Static dossier available in rva-demo-portal/out"

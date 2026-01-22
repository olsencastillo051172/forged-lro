#!/bin/bash
set -e

echo "⚔️ RVA Sovereign Audit: starting notarization..."

# 1. Generar Anchor con Go (desde cmd/rva-rotate/main.go)
echo "➡️ Generating anchor with Go..."
go run cmd/rva-rotate/main.go

# 2. Correr pruebas en Go
echo "🧪 Running Go tests..."
go test ./lib/rva/... -v

# 3. Inyectar Anchor en el Portal
echo "➡️ Injecting anchor into portal..."
cd rva-demo-portal
npm ci
npm run inject-anchor

# 4. Correr pruebas en el Portal
echo "🧪 Running Portal tests..."
npm test

# 5. Construir el Portal Forense
echo "➡️ Building forensic portal..."
npm run build

# 6. Mensaje final
echo "✅ Sovereign notarization complete."
echo "📂 Static dossier available in rva-demo-portal/out"


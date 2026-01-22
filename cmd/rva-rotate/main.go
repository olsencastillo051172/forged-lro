package main

import (
	"log"
	"os"
	"time"

	// Asegúrate de que el path coincida con tu go.mod
	"rva-core/internal/policy"
)

func main() {
	// 1. Configuración de ruta y entorno
	policyPath := os.Getenv("RVA_POLICY_PATH")
	if policyPath == "" {
		policyPath = "config/rotation_policy.json"
	}

	log.Printf("[RVA-AUDIT] Starting governance engine at %s", time.Now().UTC().Format(time.RFC3339))
	log.Printf("[RVA-AUDIT] Target policy: %s", policyPath)

	// 2. Carga de la Constitución (Loader)
	pol, err := policy.LoadPolicy(policyPath)
	if err != nil {
		log.Fatalf("[AUDIT_FAIL] Critical failure during policy loading: %v", err)
	}

	// 3. Validación de Invariantes (Validator)
	// Nota: Si estamos en modo DEV, podríamos saltar ciertas reglas, 
	// pero por ahora mantenemos el rigor total.
	if err := policy.ValidateInvariants(pol); err != nil {
		log.Fatalf("[AUDIT_FAIL] Constitution violation detected: %v", err)
	}

	// 4. Veredicto Final
	log.Println("--------------------------------------------------")
	log.Printf("VERDICT: [ALLOW_ROTATION]")
	log.Printf("ISSUER: %s (%s)", pol.Issuer.Name, pol.Issuer.ID)
	log.Printf("EPOCH_CONFIG: Interval %ds | Format: %s", 
		pol.Epochs.IntervalSeconds, 
		pol.Epochs.IDFormat,
	)
	log.Printf("SECURITY: Domain Separator [%s] is ACTIVE", pol.Constraints.DomainSeparator)
	log.Println("--------------------------------------------------")
	
	log.Println("[RVA-AUDIT] Governance check completed successfully. System is irrefutable.")
}

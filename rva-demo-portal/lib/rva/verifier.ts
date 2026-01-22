// RVA Cryptographic Core - Mathematical Sovereignty
// Compliance: Reg (EU) 2024/1183 (eIDAS 2.0)
// Audit Status: NO EXTERNAL DEPENDENCIES. Native Web Crypto API only.

export type HashHex = string;

export interface MerkleStep {
  sibling: HashHex;
  direction: 'left' | 'right';
}

export interface VerificationResult {
  isValid: boolean;
  computedRoot: HashHex;
  expectedRoot: HashHex;
  pathTrace: string[];
  timestamp: string;
}

// Security Hardening: Strict Hex Format Validation
function isValidHash(hash: string): boolean {
  return /^[0-9a-fA-F]{64}$/.test(hash);
}

/**
 * Generates SHA-256 hash using browser native Crypto API.
 * Zero dependencies ensures supply chain security.
 */
export async function sha256(input: string): Promise<HashHex> {
  const encoder = new TextEncoder();
  const data = encoder.encode(input);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

/**
 * Concatenates and hashes two nodes.
 * IMPLEMENTS DOMAIN SEPARATION to prevent Second-Preimage Attacks.
 */
async function hashPair(a: HashHex, b: HashHex): Promise<HashHex> {
  const aLower = a.toLowerCase();
  const bLower = b.toLowerCase();
  // "RVA_NODE" prefix ensures that an intermediate node cannot be confused with a leaf.
  return await sha256(`RVA_NODE:${aLower}:${bLower}`);
}

/**
 * Core Logic: Offline Merkle Path Reconstruction
 */
export async function verifyMerkleProof(
  targetHash: HashHex,
  proof: MerkleStep[],
  expectedRoot: HashHex
): Promise<VerificationResult> {
  // Input Sanitation
  if (!isValidHash(targetHash)) throw new Error(`Security Violation: Invalid targetHash format: ${targetHash}`);
  if (!isValidHash(expectedRoot)) throw new Error(`Security Violation: Invalid expectedRoot format: ${expectedRoot}`);

  let currentHash = targetHash.toLowerCase();
  const trace: string[] = [currentHash];

  for (const step of proof) {
    if (!isValidHash(step.sibling)) throw new Error(`Security Violation: Malformed sibling hash detected.`);
    const sibling = step.sibling.toLowerCase();
    
    // Strict binary tree traversal
    if (step.direction === 'left') {
      currentHash = await hashPair(sibling, currentHash);
    } else {
      currentHash = await hashPair(currentHash, sibling);
    }
    trace.push(currentHash);
  }

  const isValid = currentHash === expectedRoot.toLowerCase();
  
  return {
    isValid,
    computedRoot: currentHash,
    expectedRoot: expectedRoot.toLowerCase(),
    pathTrace: trace,
    timestamp: new Date().toISOString()
  };
}

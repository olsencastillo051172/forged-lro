// RVA Cryptographic Core — Mathematical Sovereignty
// Compliance: Reg (EU) 2024/1183 (eIDAS 2.0)
// Audit Status: HARDENED — Versioned Domain Separation & Strict Validation.

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
  timeSource: string;
}

/**
 * Strict Hexadecimal Validation
 * Checks for: null/undefined, exact length (64), and hex characters.
 */
function isValidHash(hash: unknown): hash is HashHex {
  if (typeof hash !== 'string') return false;
  if (hash.length !== 64) return false;
  return /^[0-9a-fA-F]{64}$/.test(hash);
}

/**
 * Native SHA-256 implementation via Web Crypto API.
 */
export async function sha256(input: string): Promise<HashHex> {
  const encoder = new TextEncoder();
  const data = encoder.encode(input);
  const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

/**
 * Versioned Domain Separation (v1) to prevent pre-image attacks.
 */
async function hashPair(a: HashHex, b: HashHex): Promise<HashHex> {
  const aLower = a.toLowerCase();
  const bLower = b.toLowerCase();
  // RVA_NODE:v1 ensures cross-version protocol safety.
  return await sha256(`RVA_NODE:v1:${aLower}:${bLower}`);
}

/**
 * Core Verification Logic with Enhanced Audit Metadata.
 */
export async function verifyMerkleProof(
  targetHash: HashHex,
  proof: MerkleStep[],
  expectedRoot: HashHex
): Promise<VerificationResult> {
  // 1) Inputs Check
  if (!isValidHash(targetHash)) {
    throw new Error('CRITICAL_SECURITY_ERROR: Malformed targetHash.');
  }
  if (!isValidHash(expectedRoot)) {
    throw new Error('CRITICAL_SECURITY_ERROR: Malformed expectedRoot.');
  }
  if (!Array.isArray(proof)) {
    throw new Error('CRITICAL_SECURITY_ERROR: Proof must be an array.');
  }

  let currentHash = targetHash.toLowerCase();
  const trace: string[] = [currentHash];

  // 2) Proof Traversal
  for (const [index, step] of proof.entries()) {
    if (!isValidHash(step.sibling)) {
      throw new Error(`CRITICAL_SECURITY_ERROR: Malformed sibling at proof index ${index}.`);
    }
    if (step.direction !== 'left' && step.direction !== 'right') {
      throw new Error(`CRITICAL_SECURITY_ERROR: Invalid direction at proof index ${index}.`);
    }

    const sibling = step.sibling.toLowerCase();
    currentHash =
      step.direction === 'left'
        ? await hashPair(sibling, currentHash)
        : await hashPair(currentHash, sibling);

    trace.push(currentHash);
  }

  return {
    isValid: currentHash === expectedRoot.toLowerCase(),
    computedRoot: currentHash,
    expectedRoot: expectedRoot.toLowerCase(),
    pathTrace: trace,
    timestamp: new Date().toISOString(),
    timeSource: 'Local System Clock (ISO 8601)',
  };
}

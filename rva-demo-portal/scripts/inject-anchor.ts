# File: rva-demo-portal/scripts/inject-anchor.ts

```ts
#!/usr/bin/env ts-node

/**
 * RVA Anchor Injector — Emergency Interop
 * Reads a Merkle anchor from Go output and injects it into the Next.js portal
 * so the DemoClient loads a real, end-to-end proof without manual copy/paste.
 */

import { promises as fs } from 'fs';
import * as path from 'path';

type Step = { sibling: string; direction: 'left' | 'right' };
type Anchor = {
  version: string;
  domain_separator: string;
  root: string;
  leaf: string;
  proof: Step[];
  metadata: { created_at: string; issuer: string; epoch_id: number };
};

async function main() {
  const anchorPath = process.env.RVA_ANCHOR_PATH ?? path.join('internal', 'merkle', 'anchor.json');
  const portalDataPath = path.join('rva-demo-portal', 'public', 'anchor.json');

  const raw = await fs.readFile(anchorPath, 'utf8');
  const anchor: Anchor = JSON.parse(raw);

  // Minimal sanity checks
  if (anchor.domain_separator !== 'RVA_NODE:v1') {
    throw new Error(`AUDIT_FAIL: domain_separator '${anchor.domain_separator}' is not RVA_NODE:v1`);
  }
  if (!/^[0-9a-fA-F]{64}$/.test(anchor.root) || !/^[0-9a-fA-F]{64}$/.test(anchor.leaf)) {
    throw new Error('AUDIT_FAIL: malformed root or leaf (expected 64-hex)');
  }
  for (const [i, s] of anchor.proof.entries()) {
    if (!/^[0-9a-fA-F]{64}$/.test(s.sibling)) {
      throw new Error(`AUDIT_FAIL: malformed sibling at proof index ${i}`);
    }
    if (s.direction !== 'left' && s.direction !== 'right') {
      throw new Error(`AUDIT_FAIL: invalid direction at proof index ${i}`);
    }
  }

  // Shape expected by DemoClient
  const portalPayload = {
    target: anchor.leaf.toLowerCase(),
    expectedRoot: anchor.root.toLowerCase(),
    proof: anchor.proof.map(s => ({
      sibling: s.sibling.toLowerCase(),
      direction: s.direction,
    })),
    meta: anchor.metadata,
    domain: anchor.domain_separator,
  };

  await fs.mkdir(path.dirname(portalDataPath), { recursive: true });
  await fs.writeFile(portalDataPath, JSON.stringify(portalPayload, null, 2), 'utf8');

  console.log(`[RVA-AUDIT] Anchor injected → ${portalDataPath}`);
  console.log(`Issuer: ${anchor.metadata.issuer} | Epoch: ${anchor.metadata.epoch_id}`);
  console.log(`Root: ${anchor.root}`);
}

main().catch(err => {
  console.error(`[AUDIT_FAIL] Anchor injection error: ${err.message}`);
  process.exit(1);
});

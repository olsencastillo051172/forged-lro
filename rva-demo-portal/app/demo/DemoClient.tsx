'use client';

import React, { useState, useEffect } from 'react';
// Relative import to respect the folder structure
import { verifyMerkleProof, VerificationResult } from '../../lib/rva/verifier';

type Attestation = {
  document_hash: string;
  merkle_proof: Array<{ sibling: string; direction: 'left' | 'right' }>;
  expected_root: string;
  meta: { issuer: string; epoch_id: string };
};

export default function DemoClient() {
  const [airgap, setAirgap] = useState(false);
  const [netCount, setNetCount] = useState(0);
  const [jsonInput, setJsonInput] = useState<string>('');
  const [result, setResult] = useState<VerificationResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  // --- FORENSIC INTERCEPTOR: NETWORK ISOLATION ENFORCEMENT ---
  useEffect(() => {
    // 1. Capture original browser methods before modification
    const originalFetch = window.fetch;
    const OriginalXHR = window.XMLHttpRequest;
    const originalBeacon = navigator.sendBeacon;

    // 2. Centralized Interception Logic
    const intercept = () => {
      if (airgap) {
        console.warn('RVA SECURITY AUDIT: Network attempt BLOCKED by Air-Gap Policy.');
        throw new Error('SECURITY_EXCEPTION: NETWORK_ACCESS_DENIED_BY_AIRGAP_POLICY');
      }
      setNetCount(prev => prev + 1);
    };

    // 3. Override Fetch API
    window.fetch = async (...args) => {
      intercept();
      return originalFetch(...args);
    };

    // 4. Override XMLHttpRequest (XHR)
    class PatchedXHR extends OriginalXHR {
      open(...args: any[]) { intercept(); super.open(...args); }
      send(...args: any[]) { 
        if (airgap) throw new Error('SECURITY_EXCEPTION: NETWORK_ACCESS_DENIED'); 
        super.send(...args); 
      }
    }
    // @ts-ignore
    window.XMLHttpRequest = PatchedXHR;

    // 5. Override Navigator.sendBeacon (Telemetry blocker)
    if (originalBeacon) {
      navigator.sendBeacon = (...args) => {
        intercept();
        return originalBeacon.call(navigator, ...args);
      };
    }

    // Cleanup: Restore original methods when component unmounts
    return () => {
      window.fetch = originalFetch;
      window.XMLHttpRequest = OriginalXHR;
      navigator.sendBeacon = originalBeacon;
    };
  }, [airgap]);

  const handleVerify = async () => {
    setLoading(true); setErrorMsg(null); setResult(null);
    try {
      if (!jsonInput) return;
      
      // Attempt to parse JSON
      let data: Attestation;
      try {
        data = JSON.parse(jsonInput);
      } catch (e) {
        throw new Error("Invalid Input: Please paste a valid JSON Attestation object.");
      }

      // Artificial delay (12ms) for UX processing feedback
      await new Promise(r => setTimeout(r, 12));
      
      const verification = await verifyMerkleProof(
        data.document_hash, 
        data.merkle_proof, 
        data.expected_root
      );
      setResult(verification);
    } catch (e: any) {
      console.error("Verification Error:", e);
      setErrorMsg(e.message || "Unknown Integrity Error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-black text-gray-300 font-mono p-8 max-w-4xl mx-auto">
      {/* HEADER SECTION */}
      <header className="flex justify-between items-center border-b border-gray-800 pb-4 mb-8">
        <div>
          <h1 className="text-2xl text-white font-bold tracking-tighter">RVA <span className="text-gray-600">OFFLINE</span></h1>
          <p className="text-xs text-gray-500">Sovereign Integrity Verification System</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right text-xs">
            <div className={netCount === 0 ? "text-green-400" : "text-red-400"}>NET REQ: {netCount}</div>
            <div>STATUS: {airgap ? "SECURE (OFFLINE)" : "UNSECURED (ONLINE)"}</div>
          </div>
          <button 
            onClick={() => setAirgap(!airgap)}
            style={{backgroundColor: airgap ? '#22c55e' : '#374151', color: airgap ? '#000' : '#fff'}}
            className="px-4 py-2 font-bold transition-colors uppercase text-sm"
          >
            {airgap ? 'Air-Gap: Enabled' : 'Enable Air-Gap'}
          </button>
        </div>
      </header>

      {/* ERROR FEEDBACK */}
      {errorMsg && (
        <div style={{background: 'rgba(127, 29, 29, 0.3)', borderColor: '#991b1b'}} className="border text-red-400 p-4 mb-4 text-sm font-bold">
          ⚠️ {errorMsg}
        </div>
      )}

      {/* INPUT SECTION */}
      <div className="mb-8">
        <label className="block text-xs text-gray-500 mb-2">EVIDENCE PAYLOAD (JSON)</label>
        <textarea 
          value={jsonInput}
          onChange={(e) => setJsonInput(e.target.value)}
          placeholder='Paste Attestation JSON here: {"document_hash": "...", "merkle_proof": [...]}'
          className="w-full h-40 bg-gray-900 border border-gray-800 p-4 text-xs text-white focus:outline-none focus:border-green-500 transition-colors font-mono"
        />
        <button 
          onClick={handleVerify}
          disabled={loading || !jsonInput}
          className="w-full mt-4 bg-white text-black font-bold py-3 hover:bg-gray-200 disabled:opacity-50 transition-colors uppercase tracking-wide"
        >
          {loading ? 'Verifying Mathematical Proof...' : 'Execute Integrity Check'}
        </button>
      </div>

      {/* RESULTS SECTION */}
      {result && (
        <div style={{borderColor: result.isValid ? '#22c55e' : '#ef4444'}} className="p-6 border-2 bg-gray-900/50">
          <h2 style={{color: result.isValid ? '#22c55e' : '#ef4444'}} className="text-xl font-bold mb-4">
            {result.isValid ? '✅ INTEGRITY VERIFIED (MATHEMATICALLY SOUND)' : '⛔ VERIFICATION FAILED (EVIDENCE TAMPERED)'}
          </h2>
          <div className="space-y-3 text-xs text-gray-400 mt-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="bg-gray-950 p-3 border border-gray-800">
                <strong className="text-gray-500 block mb-1">EXPECTED ROOT (ANCHOR)</strong>
                <span className="text-green-500 break-all">{result.expectedRoot}</span>
              </div>
              <div className="bg-gray-950 p-3 border border-gray-800">
                <strong className="text-gray-500 block mb-1">COMPUTED ROOT (EVIDENCE)</strong>
                <span className={`break-all ${result.isValid ? 'text-gray-300' : 'text-red-500'}`}>{result.computedRoot}</span>
              </div>
            </div>
            
            <div className="mt-4 pt-4 border-t border-gray-800">
               <strong className="text-gray-500">AUDIT TRAIL:</strong>
               <p className="mt-1">Timestamp: {result.timestamp}</p>
               <p>Network Requests during session: {netCount}</p>
               <p>Environment: {airgap ? "Air-Gapped (Compliant)" : "Networked (Non-Compliant)"}</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

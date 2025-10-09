# Threat Model – Bridge (STRIDE)

## Scope

- Message flow: Send → Verify → Execute  
- Clients: light-client, zk-client  
- Risk Layer: veto, rate-limits, circuit breakers  
- Key assets: funds, headers, proofs, message store  

## STRIDE

### S – Spoofing

- **Mitigation:** strong identity for relayers, proof-bound verification, client auth.  

### T – Tampering

- **Mitigation:** proofs validate headers; immutable ID; signed releases; SBOM.  

### R – Repudiation

- **Mitigation:** on-chain events; audit trail; signed governance proposals.  

### I – Information Disclosure

- **Mitigation:** minimize sensitive data on-chain; redact PII; encrypt ops configs.  

### D – Denial of Service

- **Mitigation:** rate-limits per asset/route/account; circuit breakers; mempool guards.  

### E – Elevation of Privilege

- **Mitigation:** time-locked upgrades; multi-sig; separation of duties; no hidden bypass.  

## Invariants

- No Execute without successful Verify.  
- One-time execution per (ID, Nonce).  
- Client diversity: failure in one client does not approve execution.  
- Upgrades cannot disable verification or replay protections.  

## Tests

- Property-based tests on idempotency and invariants.  
- Fuzzing malformed proofs/headers.  
- Replay attempts; upgrade attempts; kill-switch exercises.  

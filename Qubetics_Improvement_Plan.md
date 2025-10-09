# RFC-001: Qubetics Improvement Plan (Bridge + Security + Ops)

## Motivation
Price collapse and trust issues were driven by linear unlock pressure and weak engineering process. We propose a phased, governance-driven technical upgrade: verifiable bridges (light/zk-clients), risk layer, auditable on-chain vesting, and production SDLC.

## Scope (Phases)
- **Phase A (Weeks 1–4):** On-chain vesting `x/tics`; CI/CD; security hardening; public telemetry and daily snapshots.
- **Phase B (Weeks 5–12):** Bridge `x/bridge` MVP with light-clients (ETH/SOL) on testnet; Risk Layer design; rate limits & circuit breakers.
- **Phase C (Weeks 13–20):** zk-light-clients integration; Risk Layer go-live; audits (2 vendors), formal verification for invariants.
- **Phase D (Governance Upgrade):** mainnet roll-out via gov proposal; canary release; post-mortem drills.

## Deliverables
- `x/tics` (vesting): TGE + daily %, claim logic, events, invariants, queries.
- `x/bridge` (MVP): Send → Verify (light/zk proof) → Execute; idempotency; replay protection.
- Risk Layer: independent veto for anomalies; rate-limits; kill-switch; time-locked upgrades.
- SDLC: CI (build/test/lint/gosec), CODEOWNERS, SECURITY.md, CONTRIBUTING, issue/PR templates.
- Ops: systemd unit, Ansible role, snapshots exporter, Prometheus/Grafana dashboards.

## Success Criteria
- No manual token transfers: all unlocks/claims on-chain and auditable.
- Bridge ops require valid proofs; no execute without verify; dual-client checks.
- Risk Layer can pause/limit flows during anomalies.
- 99.9% availability targets with documented SLO and recovery runbooks.

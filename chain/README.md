# Qubetics Chain Add-ons

This folder contains community add-ons for the Qubetics stack:

- `x/bridge/` — prototype of cross-chain bridge module (Send → Verify → Execute).
- `x/tics/` (reserved) — vesting/claims.
- Shared CI: unit tests + gosec; SBOM + OSS scanning in the repo root.

## Requirements

- Go **1.23.x**
- (for module dev) Cosmos SDK **0.47.x** line
- `make` (optional)

## Quick start

```bash
cd chain/x/bridge
go mod tidy
go test ./... -v

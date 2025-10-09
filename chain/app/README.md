# QubeticsApp Integration Layer

This is the main application entrypoint for the modular Qubetics chain.

## Modules Integrated

- **x/bridge** — cross-chain messaging, proof validation, execution.
- **x/tics** — vesting, airdrop and on-chain claim logic.
- **x/governance** — DAO-based control with timelock and multisig.
- **Risk Layer** — rate limits, kill-switch, audit hooks.

## Notes

This app skeleton demonstrates how Cosmos SDK-based modular architecture can be wired together.

Further steps:

- Register message routes.
- Add CLI and gRPC interfaces.
- Integrate consensus engine (CometBFT/Tendermint).

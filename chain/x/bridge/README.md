# Bridge Keeper — Overview

The keeper manages message state transitions within the bridge:
- **InitGenesis()** — module bootstrap.
- **HandleMessage()** — receives and routes bridge messages.
- **Hooks:** rate-limit checks, proof validation, and event emission.

### Message Flow

`Send → Verify → Execute`
Each step emits standardized events:
- `bridge_send`
- `bridge_verify`
- `bridge_execute`

The keeper ensures **idempotency**, **proof validation**, and **replay protection**.

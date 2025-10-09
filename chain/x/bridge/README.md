# x/bridge (skeleton)

Principles:
- Two-phase flow: **Send -> Verify (proof) -> Execute**
- **No Execute without successful Verify**
- **Idempotent** Execute (replays are dropped)
- **Dual client support** (e.g., light client + zk client) with quorum policy
- **Risk layer hooks**: pause/limits when anomalies detected

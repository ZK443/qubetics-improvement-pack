# Qubetics Bridge (skeleton, "right" variant)

**Flow:** `Send → Verify → Execute`

## Инварианты

- **Нет Execute без Verify** (подтверждённого доказ-ва)
- **Идемпотентность:** повторный Execute для одного `msg.ID` — no-op
- **Replay protection:** монотонный `Nonce` per-route
- **Risk-layer:** rate-limits, pause/kill-switch, timelock на апгрейды (далее)

## Поток

1) **Send** — создаётся заявка (`MsgSend`), пишется `StatusPending`.
2) **Verify** — хранится валидный `Proof`, статус `StatusVerified`.
3) **Execute** — проверка статуса, применение эффекта, `StatusExecuted`.

## Клиенты доказательств

- **ETH**: light+Merkle (опц. zk)
- **SOL**: light+Merkle
- **BTC**: SPV (headers+tx) / HTLC

## Хранилище (минимум)

- `msg:ID`
- `st:ID`
- `pf:ID`
- `rl:route`
- `ps:route`

## Метрики/аудит

- События: `bridge_send`, `bridge_verify`, `bridge_execute`
- Дашборды: `docs/grafana/*`

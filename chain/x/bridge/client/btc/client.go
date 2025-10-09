// SPDX-License-Identifier: MIT
package btc

import qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// ProofBTC — SPV-пакет: цепочка заголовков + merkle-ветка транзакции.
type ProofBTC struct {
	Headers    [][]byte // последовательность хедеров (для достаточной работы)
	TxID       []byte   // txid (little-endian или big-endian — договоримся)
	MerklePath [][]byte // ветка к txid
}

type Client struct{}

func (Client) Network() string { return "bitcoin" }

func (Client) Verify(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult {
	return qtypes.VerificationResult{Valid: false, Reason: "btc spv verification not implemented"}
}

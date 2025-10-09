// SPDX-License-Identifier: MIT
package sol

import qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// ProofSOL — заглушка доказательства для Solana (headers + merkle).
type ProofSOL struct {
	BlockHeader []byte
	MerklePath  []byte
	Index       uint32
}

type Client struct{}

func (Client) Network() string { return "solana" }

func (Client) Verify(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult {
	return qtypes.VerificationResult{Valid: false, Reason: "sol light verification not implemented"}
}

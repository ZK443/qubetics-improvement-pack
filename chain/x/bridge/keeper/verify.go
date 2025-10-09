// SPDX-License-Identifier: MIT
package keeper

import qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

func (k Keeper) VerifyWithLightClient(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult {
	return qtypes.VerificationResult{Valid: false, Reason: "light client not implemented"}
}
func (k Keeper) VerifyWithZK(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult {
	return qtypes.VerificationResult{Valid: false, Reason: "zk verification not implemented"}
}
func (k Keeper) VerifyWithSPV(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult {
	return qtypes.VerificationResult{Valid: false, Reason: "SPV not implemented"}
}

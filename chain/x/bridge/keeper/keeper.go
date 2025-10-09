package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"qubetics/chain/x/bridge/types"
)

type Keeper struct {
	bank    types.BankKeeper // interface to send coins if needed
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, bank types.BankKeeper) Keeper {
	return Keeper{storeKey: storeKey, bank: bank}
}

// VerifyProof must be pure: no state mutation on verification only.
func (k Keeper) VerifyProof(ctx sdk.Context, msg types.Message, proof types.Proof) types.VerificationResult {
	// TODO: call registered client(s) to verify proof (light/zk).
	// MUST NOT execute state changes here.
	return types.VerificationResult{Valid: false, Reason: "not implemented"}
}

// Execute applies the effect of a previously verified message.
// MUST be idempotent and check replay by msg.ID/Nonce.
func (k Keeper) Execute(ctx sdk.Context, msg types.Message) error {
	// TODO: idempotency check, route dispatch, replay protection.
	return nil
}

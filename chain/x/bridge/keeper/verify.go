package keeper

import (
    "context"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func (k Keeper) Verify(ctx context.Context, msg types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
    store := sdk.UnwrapSDKContext(ctx).KVStore(k.storeKey)
    proofKey := []byte(msg.ProofId)

    if store.Has(proofKey) {
        return nil, types.ErrProofAlreadyExists
    }

    // TODO: Plug in light/zk proof validation here
    store.Set(proofKey, []byte(msg.ProofData))

    ctxSDK := sdk.UnwrapSDKContext(ctx)
    ctxSDK.EventManager().EmitEvent(
        sdk.NewEvent("bridge_verify",
            sdk.NewAttribute("proof_id", msg.ProofId),
            sdk.NewAttribute("verifier", msg.Verifier),
        ),
    )

    return &types.MsgVerifyProofResponse{Status: "verified"}, nil
}

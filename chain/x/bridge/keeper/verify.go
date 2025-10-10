package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func (k Keeper) Verify(ctx context.Context, msg types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)

    // 1) Авторизация
    if !k.IsRelayer(sdkCtx, msg.Signer) {
        return nil, types.ErrUnauthorized
    }
    if len(msg.ProofData) == 0 || msg.ProofId == "" || msg.MessageID == "" {
        return nil, types.ErrInvalidProof
    }

    // 2) Биндинг (proof должен покрывать binding(messageID, expectedDigest))
    binding := ComputeBindingHash(msg.MessageID, msg.ExpectedDigest)
    if k.bindingVerifier != nil && !k.bindingVerifier.VerifyBinding(binding, msg.ProofData) {
        return nil, types.ErrInvalidProof
    }

    // 3) Идемпотентность: запрещаем повторную верификацию
    store := sdkCtx.KVStore(k.storeKey)
    pKey := k.proofKey(msg.ProofId)
    if store.Has(pKey) {
        return nil, types.ErrProofAlreadyExists
    }
    store.Set(pKey, []byte("verified"))

    // 4) Событие
    sdkCtx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_verify",
            sdk.NewAttribute("proof_id", msg.ProofId),
            sdk.NewAttribute("verifier", msg.Verifier),
            sdk.NewAttribute("message_id", msg.MessageID),
        ),
    )

    return &types.MsgVerifyProofResponse{Status: "verified"}, nil
}

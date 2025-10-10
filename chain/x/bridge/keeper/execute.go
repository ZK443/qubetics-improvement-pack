package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func (k Keeper) Execute(ctx context.Context, msg types.MsgExecute) (*types.MsgExecuteResponse, error) {
    sdkCtx := sdk.UnwrapSDKContext(ctx)

    // Pause / Rate-limits (если реализованы)
    if k.IsBridgePaused(sdkCtx) {
        return nil, types.ErrBridgePaused
    }
    if !k.CheckRateLimit(sdkCtx, msg.Executor, msg.Amount) {
        return nil, types.ErrRateLimitExceeded
    }

    store := sdkCtx.KVStore(k.storeKey)
    eKey := k.execKey(msg.MessageId)
    pKey := k.proofKey(msg.ProofId)

    // Проверяем, что Verify был успешен
    if !store.Has(pKey) {
        return nil, types.ErrNotVerified
    }

    // Идемпотентность: атомарная проверка + запись статуса исполнения
    if store.Has(eKey) {
        return nil, types.ErrAlreadyExecuted
    }
    store.Set(eKey, []byte("executed"))

    // TODO: здесь должен быть реальный mint/unlock/contract-call, зависит от интеграции

    sdkCtx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_execute",
            sdk.NewAttribute("message_id", msg.MessageId),
            sdk.NewAttribute("executor", msg.Executor.String()),
        ),
    )

    return &types.MsgExecuteResponse{Status: "done"}, nil
}

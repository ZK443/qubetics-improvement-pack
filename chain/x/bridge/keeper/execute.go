package keeper

import (
    "context"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func (k Keeper) Execute(ctx context.Context, msg types.MsgExecute) (*types.MsgExecuteResponse, error) {
    store := sdk.UnwrapSDKContext(ctx).KVStore(k.storeKey)
    execKey := []byte(msg.MessageId)

    if store.Has(execKey) {
        return nil, types.ErrAlreadyExecuted
    }

    // TODO: Execute token mint/unlock logic based on message content
    store.Set(execKey, []byte("executed"))

    sdkCtx := sdk.UnwrapSDKContext(ctx)
    sdkCtx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_execute",
            sdk.NewAttribute("message_id", msg.MessageId),
            sdk.NewAttribute("executor", msg.Executor),
        ),
    )

    return &types.MsgExecuteResponse{Status: "done"}, nil
}

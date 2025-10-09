package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type Approval struct {
    Action string
    Approver sdk.AccAddress
}

func (k Keeper) ApproveAction(ctx sdk.Context, action string, approver sdk.AccAddress) {
    store := ctx.KVStore(k.storeKey)
    key := []byte("approval:" + action + ":" + approver.String())
    store.Set(key, []byte("ok"))

    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_action_approved",
            sdk.NewAttribute("action", action),
            sdk.NewAttribute("approver", approver.String()),
        ),
    )
}

func (k Keeper) IsActionApproved(ctx sdk.Context, action string, threshold int) bool {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, []byte("approval:"+action))
    defer iterator.Close()

    count := 0
    for ; iterator.Valid(); iterator.Next() {
        count++
        if count >= threshold {
            return true
        }
    }
    return false
}

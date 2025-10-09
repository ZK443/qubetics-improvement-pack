package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) KillSwitch(ctx sdk.Context, reason string) {
    store := ctx.KVStore(k.storeKey)
    store.Set([]byte("kill_switch"), []byte(reason))

    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_paused",
            sdk.NewAttribute("reason", reason),
        ),
    )
}

func (k Keeper) IsBridgePaused(ctx sdk.Context) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has([]byte("kill_switch"))
}

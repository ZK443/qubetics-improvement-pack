package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsAdmin(ctx sdk.Context, addr sdk.AccAddress) bool {
    store := ctx.KVStore(k.storeKey)
    key := []byte("admin:" + addr.String())
    return store.Has(key)
}

func (k Keeper) AddAdmin(ctx sdk.Context, addr sdk.AccAddress) {
    store := ctx.KVStore(k.storeKey)
    key := []byte("admin:" + addr.String())
    store.Set(key, []byte("true"))

    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_admin_added",
            sdk.NewAttribute("admin", addr.String()),
        ),
    )
}

func (k Keeper) RemoveAdmin(ctx sdk.Context, addr sdk.AccAddress) {
    store := ctx.KVStore(k.storeKey)
    key := []byte("admin:" + addr.String())
    store.Delete(key)

    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_admin_removed",
            sdk.NewAttribute("admin", addr.String()),
        ),
    )
}

package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "time"
)

var (
    DefaultMaxPerAccount = sdk.NewInt(10_000_000_000) // пример: 10 QUBT за час
    DefaultCooldown      = time.Hour
)

func (k Keeper) CheckRateLimit(ctx sdk.Context, account sdk.AccAddress, amount sdk.Int) bool {
    store := ctx.KVStore(k.storeKey)
    key := []byte("limit:" + account.String())

    bz := store.Get(key)
    if bz != nil {
        last := sdk.BigEndianToUint64(bz)
        if ctx.BlockTime().Sub(time.Unix(int64(last), 0)) < DefaultCooldown {
            return false
        }
    }
    if amount.GT(DefaultMaxPerAccount) {
        return false
    }
    store.Set(key, sdk.Uint64ToBigEndian(uint64(ctx.BlockTime().Unix())))
    return true
}

package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "time"
)

type PendingAction struct {
    Action string
    ETA    time.Time
    Proposer sdk.AccAddress
}

func (k Keeper) ProposeUpgrade(ctx sdk.Context, proposer sdk.AccAddress, action string, delay time.Duration) {
    store := ctx.KVStore(k.storeKey)
    eta := ctx.BlockTime().Add(delay)

    pa := PendingAction{
        Action:   action,
        ETA:      eta,
        Proposer: proposer,
    }

    store.Set([]byte("proposal:"+action), []byte(eta.Format(time.RFC3339)))
    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_upgrade_proposed",
            sdk.NewAttribute("action", action),
            sdk.NewAttribute("eta", eta.String()),
        ),
    )
}

func (k Keeper) ExecuteUpgrade(ctx sdk.Context, action string) bool {
    store := ctx.KVStore(k.storeKey)
    etaBytes := store.Get([]byte("proposal:" + action))
    if etaBytes == nil {
        return false
    }

    eta, _ := time.Parse(time.RFC3339, string(etaBytes))
    if ctx.BlockTime().Before(eta) {
        return false // still timelocked
    }

    store.Delete([]byte("proposal:" + action))
    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_upgrade_executed",
            sdk.NewAttribute("action", action),
        ),
    )
    return true
}

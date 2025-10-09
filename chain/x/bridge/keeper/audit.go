package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "fmt"
)

func (k Keeper) EmitAuditLog(ctx sdk.Context, action, actor, details string) {
    ctx.EventManager().EmitEvent(
        sdk.NewEvent("bridge_audit",
            sdk.NewAttribute("action", action),
            sdk.NewAttribute("actor", actor),
            sdk.NewAttribute("details", fmt.Sprintf("%.100s", details)),
        ),
    )
}

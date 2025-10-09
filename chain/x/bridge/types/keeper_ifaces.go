package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper is a minimal example interface (extend as needed).
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

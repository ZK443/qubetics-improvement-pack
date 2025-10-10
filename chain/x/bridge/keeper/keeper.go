package keeper

import (
    "github.com/cosmos/cosmos-sdk/store/prefix"
    storetypes "github.com/cosmos/cosmos-sdk/store/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

type Keeper struct {
    storeKey       storetypes.StoreKey
    bindingVerifier types.BindingVerifier
}

func NewKeeper(key storetypes.StoreKey, bv types.BindingVerifier) Keeper {
    return Keeper{
        storeKey:       key,
        bindingVerifier: bv,
    }
}

// Helpers
func (k Keeper) proofKey(id string) []byte { return []byte(types.KeyPrefixProof + id) }
func (k Keeper) execKey(id string)  []byte { return []byte(types.KeyPrefixExec + id) }

func (k Keeper) IsRelayer(ctx sdk.Context, addr sdk.AccAddress) bool {
    // Упростим: любой не-пустой адрес считаем валидным для PoC, в реале — роли/ACL.
    if addr == nil || addr.Empty() { return false }
    return true
}

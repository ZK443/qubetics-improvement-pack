//go:build cosmos

package keeper

import (
	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Cosmos-style Keeper: KV-хранилище, codec, params.
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	ps       paramtypes.Subspace
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ps paramtypes.Subspace) Keeper {
	return Keeper{cdc: cdc, storeKey: key, ps: ps}
}

// Ниже — минимальные заглушки под те же методы, что у in-memory,
// чтобы код, который их вызывает, собирался и под тегом cosmos.
func (k Keeper) isPaused(ctx sdk.Context) bool { // пример ctx-варианта
	// TODO: читать флаг из KV: store.Get(KeyPause()) == []byte{1}
	return false
}
func (k Keeper) getStatusByID(ctx sdk.Context, id string) types.Status {
	// TODO: decode из KV по KeyStatus(id)
	return types.StatusUnknown
}
func (k Keeper) isExecuted(ctx sdk.Context, id string) bool {
	// TODO: по KeyStatus(id) == StatusExecuted
	return false
}
func (k Keeper) rateLimited(_ sdk.Context, _ types.MsgExecute) (bool, string) {
	// TODO: проверка квот по KeyRateLimit(route)
	return false, ""
}
func (k Keeper) markExecuted(ctx sdk.Context, id string) {
	// TODO: store.Set(KeyStatus(id), StatusExecuted)
	_ = id
}
func (k Keeper) emitEvent(_ sdk.Context, _ string, _ map[string]string) {}

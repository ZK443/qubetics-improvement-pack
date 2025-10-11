//go:build cosmos

package keeper

import (
	"encoding/json"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/store"
	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	ps       paramtypes.Subspace
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ps paramtypes.Subspace) Keeper {
	return Keeper{cdc: cdc, storeKey: key, ps: ps}
}

// ---- helpers ----
func (k Keeper) kv(ctx sdk.Context) sdk.KVStore { return ctx.KVStore(k.storeKey) }

// ---- Params ----
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	kv := k.kv(ctx)
	bz := kv.Get(store.KeyParams)
	if bz == nil {
		return types.DefaultParams()
	}
	var p types.Params
	if err := json.Unmarshal(bz, &p); err != nil {
		// fallback к дефолту при повреждённых данных
		return types.DefaultParams()
	}
	return p
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}
	bz, _ := json.Marshal(p)
	k.kv(ctx).Set(store.KeyParams, bz)
	return nil
}

// ---- ACL ----
func (k Keeper) IsAllowed(ctx sdk.Context, addr string) bool {
	bz := k.kv(ctx).Get(append(store.KeyACL, []byte(addr)...))
	return len(bz) == 1 && bz[0] == 1
}

func (k Keeper) SetAllowed(ctx sdk.Context, addr string, allowed bool) {
	val := byte(0)
	if allowed {
		val = 1
	}
	k.kv(ctx).Set(append(store.KeyACL, []byte(addr)...), []byte{val})
}

// прежние заглушки (контекстные варианты); реализуются в S4.
func (k Keeper) isPaused(ctx sdk.Context) bool                                 { return k.GetParams(ctx).GlobalPause }
func (k Keeper) getStatusByID(_ sdk.Context, _ string) types.Status            { return types.StatusUnknown }
func (k Keeper) isExecuted(_ sdk.Context, _ string) bool                       { return false }
func (k Keeper) rateLimited(_ sdk.Context, _ types.MsgExecute) (bool, string)  { return false, "" }
func (k Keeper) markExecuted(_ sdk.Context, _ string)                           {}
func (k Keeper) emitEvent(_ sdk.Context, _ string, _ map[string]string)         {}

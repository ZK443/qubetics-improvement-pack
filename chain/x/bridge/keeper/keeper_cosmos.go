//go:build cosmos

package keeper

import (
	"encoding/binary"
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
		// при повреждённых данных — безопасный дефолт
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

// ---- Status CRUD ----
func (k Keeper) GetStatus(ctx sdk.Context, id string) types.Status {
	bz := k.kv(ctx).Get(append([]byte("bridge/status/"), []byte(id)...))
	if len(bz) != 1 {
		return types.StatusUnknown
	}
	return types.Status(bz[0])
}

func (k Keeper) SetStatus(ctx sdk.Context, id string, st types.Status) {
	k.kv(ctx).Set(append([]byte("bridge/status/"), []byte(id)...), []byte{byte(st)})
}

// ---- Nonce per-sender ----
func (k Keeper) PeekNonce(ctx sdk.Context, sender string) uint64 {
	bz := k.kv(ctx).Get(append([]byte("bridge/nonce/"), []byte(sender)...))
	if len(bz) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) NextNonce(ctx sdk.Context, sender string) uint64 {
	cur := k.PeekNonce(ctx, sender) + 1
	var out [8]byte
	binary.BigEndian.PutUint64(out[:], cur)
	k.kv(ctx).Set(append([]byte("bridge/nonce/"), []byte(sender)...), out[:])
	return cur
}

// ---- Invariants / Guards ----
func (k Keeper) CanExecute(ctx sdk.Context, id string) bool {
	if k.isPaused(ctx) {
		return false
	}
	if k.GetStatus(ctx, id) != types.StatusVerified {
		return false
	}
	if k.isExecuted(ctx, id) {
		return false
	}
	return true
}

func (k Keeper) MarkExecuted(ctx sdk.Context, id string) {
	k.markExecuted(ctx, id)
	k.SetStatus(ctx, id, types.StatusExecuted)
}

// Контекстные заглушки для совместимости (реализация на следующих этапах).
func (k Keeper) isPaused(ctx sdk.Context) bool                                { return k.GetParams(ctx).GlobalPause }
func (k Keeper) getStatusByID(_ sdk.Context, _ string) types.Status           { return types.StatusUnknown }
func (k Keeper) isExecuted(_ sdk.Context, _ string) bool                      { return false }
func (k Keeper) rateLimited(_ sdk.Context, _ types.MsgExecute) (bool, string) { return false, "" }
func (k Keeper) markExecuted(_ sdk.Context, _ string)                          {}
func (k Keeper) emitEvent(_ sdk.Context, _ string, _ map[string]string)        {}

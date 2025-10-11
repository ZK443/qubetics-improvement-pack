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
	key := append([]byte("bridge/status/"), []byte(id)...)
	bz := k.kv(ctx).Get(key)
	if len(bz) != 1 {
		return types.StatusUnknown
	}
	return types.Status(bz[0])
}

func (k Keeper) SetStatus(ctx sdk.Context, id string, st types.Status) {
	key := append([]byte("bridge/status/"), []byte(id)...)
	k.kv(ctx).Set(key, []byte{byte(st)})
}

// ---- Nonce per-sender ----
func (k Keeper) PeekNonce(ctx sdk.Context, sender string) uint64 {
	key := append([]byte("bridge/nonce/"), []byte(sender)...)
	bz := k.kv(ctx).Get(key)
	if len(bz) != 8 {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}
func (k Keeper) NextNonce(ctx sdk.Context, sender string) uint64 {
	key := append([]byte("bridge/nonce/"), []byte(sender)...)
	cur := k.PeekNonce(ctx, sender) + 1
	var out [8]byte
	binary.BigEndian.PutUint64(out[:], cur)
	k.kv(ctx).Set(key, out[:])
	return cur
}

// ---- Rate-limit (окна по маршруту) ----
// Логика идентична in-memory варианту, но хранится в KV:
// count: bridge/rl/count/<route>  (uint64)
// until: bridge/rl/until/<route>  (int64, unix ms)
func (k Keeper) rateLimited(ctx sdk.Context, msg types.MsgExecute) (bool, string) {
	p := k.GetParams(ctx)
	if p.RateLimitAmount == 0 || p.RateLimitWindowMs == 0 {
		return false, ""
	}

	route := string(msg.Route)
	now := ctx.BlockTime().UnixMilli()

	// load until
	kUntil := store.KeyRLUntil(route)
	untilBz := k.kv(ctx).Get(kUntil)
	var until int64
	if len(untilBz) == 8 {
		until = int64(binary.BigEndian.Uint64(untilBz))
	}

	// если окно истекло — открыть новое
	if now > until {
		until = now + int64(p.RateLimitWindowMs)
		var out [8]byte
		binary.BigEndian.PutUint64(out[:], uint64(until))
		k.kv(ctx).Set(kUntil, out[:])

		// обнулить count
		k.kv(ctx).Set(store.KeyRLCount(route), make([]byte, 8))
	}

	// load count
	cntBz := k.kv(ctx).Get(store.KeyRLCount(route))
	var cnt uint64
	if len(cntBz) == 8 {
		cnt = binary.BigEndian.Uint64(cntBz)
	}

	if cnt >= p.RateLimitAmount {
		return true, "window"
	}

	// increment count
	cnt++
	var out [8]byte
	binary.BigEndian.PutUint64(out[:], cnt)
	k.kv(ctx).Set(store.KeyRLCount(route), out[:])

	return false, ""
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

// ---- Events ----
func (k Keeper) emitEvent(ctx sdk.Context, evt string, attrs map[string]string) {
	ev := sdk.NewEvent(evt)
	for kAttr, vAttr := range attrs {
		ev = ev.AppendAttributes(sdk.NewAttribute(kAttr, vAttr))
	}
	ctx.EventManager().EmitEvent(ev)
}

// Контекстные заглушки (реализация будет дорабатываться на следующих этапах).
func (k Keeper) isPaused(ctx sdk.Context) bool                                { return k.GetParams(ctx).GlobalPause }
func (k Keeper) getStatusByID(_ sdk.Context, _ string) types.Status           { return types.StatusUnknown }
func (k Keeper) isExecuted(_ sdk.Context, _ string) bool                      { return false }
func (k Keeper) rateLimitedLegacy(_ sdk.Context, _ types.MsgExecute) (bool, string) {
	return false, ""
}
func (k Keeper) markExecuted(_ sdk.Context, _ string)                   {}
func (k Keeper) emitEventLegacy(_ sdk.Context, _ string, _ map[string]string) {}

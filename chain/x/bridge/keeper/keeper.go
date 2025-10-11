//go:build !cosmos

package keeper

import (
	"time"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// Лёгкий in-memory Keeper для прототипа и CI (без Cosmos SDK).
type Keeper struct {
	paused   bool
	executed map[string]bool
	status   map[string]types.Status
	nonce    map[string]uint64

	// rate-limit (окна по маршруту)
	rlCount map[types.Route]uint64 // сколько выполнений в текущем окне
	rlUntil map[types.Route]int64  // unix ms — когда заканчивается окно

	// источник времени (подменяется в тестах)
	now func() int64

	params types.Params
	acl    map[string]bool // "адрес" -> разрешён
}

func NewKeeper() *Keeper {
	return &Keeper{
		executed: make(map[string]bool),
		status:   make(map[string]types.Status),
		nonce:    make(map[string]uint64),
		rlCount:  map[types.Route]uint64{},
		rlUntil:  map[types.Route]int64{},
		now:      func() int64 { return time.Now().UnixMilli() },

		params: types.DefaultParams(),
		acl:    map[string]bool{},
	}
}

// ---- базовые методы ----
func (k *Keeper) isPaused() bool                                { return k.paused || k.params.GlobalPause }
func (k *Keeper) getStatusByID(id string) types.Status          { return k.status[id] }
func (k *Keeper) isExecuted(id string) bool                     { return k.executed[id] }
func (k *Keeper) markExecuted(id string)                        { k.executed[id] = true }
func (k *Keeper) emitEvent(_ string, _ map[string]string)       {}

// Минимальная реализация rate-limit поверх Params.
// Считает количество выполнений по маршруту в скользящем окне.
func (k *Keeper) rateLimited(msg types.MsgExecute) (bool, string) {
	// Защита от "выключенного" лимита (на всякий случай).
	if k.params.RateLimitAmount == 0 || k.params.RateLimitWindowMs == 0 {
		return false, ""
	}

	now := k.now()
	until := k.rlUntil[msg.Route]
	// Если окно истекло — открыть новое окно.
	if now > until {
		k.rlUntil[msg.Route] = now + int64(k.params.RateLimitWindowMs)
		k.rlCount[msg.Route] = 0
	}
	if k.rlCount[msg.Route] >= k.params.RateLimitAmount {
		return true, "window"
	}
	k.rlCount[msg.Route]++
	return false, ""
}

// ---- Params ----
func (k *Keeper) GetParams() types.Params { return k.params }

func (k *Keeper) SetParams(p types.Params) error {
	if err := p.Validate(); err != nil {
		return err
	}
	k.params = p
	return nil
}

// ---- ACL ----
func (k *Keeper) IsAllowed(addr string) bool {
	allowed, ok := k.acl[addr]
	return ok && allowed
}

func (k *Keeper) SetAllowed(addr string, allowed bool) {
	if k.acl == nil {
		k.acl = map[string]bool{}
	}
	k.acl[addr] = allowed
}

// ---- Status CRUD ----
func (k *Keeper) GetStatus(id string) types.Status {
	if st, ok := k.status[id]; ok {
		return st
	}
	return types.StatusUnknown
}

func (k *Keeper) SetStatus(id string, st types.Status) {
	k.status[id] = st
}

// ---- Nonce per-sender ----
func (k *Keeper) PeekNonce(sender string) uint64 {
	return k.nonce[sender]
}

func (k *Keeper) NextNonce(sender string) uint64 {
	n := k.nonce[sender] + 1
	k.nonce[sender] = n
	return n
}

// ---- Invariants / Guards ----
// Можно выполнять, если: нет глобальной паузы, статус == Verified, не выполнено ранее.
func (k *Keeper) CanExecute(id string) bool {
	if k.isPaused() {
		return false
	}
	if k.GetStatus(id) != types.StatusVerified {
		return false
	}
	if k.isExecuted(id) {
		return false
	}
	return true
}

// Зафиксировать успешное выполнение с обновлением статуса и флагов.
func (k *Keeper) MarkExecuted(id string) {
	k.markExecuted(id)               // внутренний быстрый флаг
	k.SetStatus(id, types.StatusExecuted)
}

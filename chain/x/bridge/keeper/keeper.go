//go:build !cosmos

package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// Лёгкий in-memory Keeper для прототипа и CI.
type Keeper struct {
	paused   bool
	executed map[string]bool
	status   map[string]types.Status
	nonce    map[string]uint64

	params types.Params
	acl    map[string]bool // "адрес" -> разрешён
}

func NewKeeper() *Keeper {
	return &Keeper{
		executed: make(map[string]bool),
		status:   make(map[string]types.Status),
		nonce:    make(map[string]uint64),
		params:   types.DefaultParams(),
		acl:      map[string]bool{},
	}
}

// ---- базовые ранее добавленные методы опущены для краткости ----

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

// Сообщение можно выполнять, если:
// 1) глобальная пауза выключена; 2) статус == Verified; 3) ранее не выполнено.
func (k *Keeper) CanExecute(id string) bool {
	if k.isPaused() {
		return false
	}
	st := k.GetStatus(id)
	if st != types.StatusVerified {
		return false
	}
	if k.isExecuted(id) {
		return false
	}
	return true
}

// Зафиксировать успешное выполнение с обновлением статуса и флагов.
func (k *Keeper) MarkExecuted(id string) {
	k.markExecuted(id)         // внутренний флаг быстрого пути
	k.SetStatus(id, types.StatusExecuted)
}

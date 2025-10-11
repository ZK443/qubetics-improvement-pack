//go:build !cosmos

package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// Лёгкий in-memory Keeper для прототипа и CI.
type Keeper struct {
	paused   bool
	executed map[string]bool
	status   map[string]types.Status

	params types.Params
	acl    map[string]bool // "адрес" -> разрешён
}

func NewKeeper() *Keeper {
	return &Keeper{
		executed: make(map[string]bool),
		status:   make(map[string]types.Status),
		params:   types.DefaultParams(),
		acl:      map[string]bool{},
	}
}

func (k *Keeper) isPaused() bool                               { return k.paused || k.params.GlobalPause }
func (k *Keeper) getStatusByID(id string) types.Status         { return k.status[id] }
func (k *Keeper) isExecuted(id string) bool                    { return k.executed[id] }
func (k *Keeper) rateLimited(_ types.MsgExecute) (bool, string) { // подробная логика — в S4
	return false, ""
}
func (k *Keeper) markExecuted(id string)                  { k.executed[id] = true }
func (k *Keeper) emitEvent(_ string, _ map[string]string) {}

// ---- Params ----
func (k *Keeper) GetParams() types.Params        { return k.params }
func (k *Keeper) SetParams(p types.Params) error { // валидация — из types.Params
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

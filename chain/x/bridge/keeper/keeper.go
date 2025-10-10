package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// Лёгкий in-memory Keeper для прототипа и тестов CI.
// Не тянем Cosmos SDK, чтобы держать сборку быстрой и простой.
type Keeper struct {
	paused   bool
	executed map[string]bool
	status   map[string]types.Status
}

func NewKeeper() *Keeper {
	return &Keeper{
		executed: make(map[string]bool),
		status:   make(map[string]types.Status),
	}
}

func (k *Keeper) isPaused() bool                            { return k.paused }
func (k *Keeper) getStatusByID(id string) types.Status      { return k.status[id] }
func (k *Keeper) isExecuted(id string) bool                 { return k.executed[id] }
func (k *Keeper) rateLimited(_ types.MsgExecute) (bool, string) { return false, "" }
func (k *Keeper) markExecuted(id string)                    { k.executed[id] = true }
func (k *Keeper) emitEvent(_ string, _ map[string]string)   {}

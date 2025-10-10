package keeper

import (
	qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// Execute — идемпотентное применение эффекта после успешного Verify.
// Прототип: возвращает статус и текстовое объяснение.
func (k Keeper) Execute(msg qtypes.Message) (qtypes.Status, string) {
	// 1) Глобальная пауза
	if k.isPaused() {
		return qtypes.StatusRejected, "paused"
	}

	// 2) Должен быть Verify
	st := k.getStatusByID(msg.ID)
	if st != qtypes.StatusVerified {
		return qtypes.StatusRejected, "not-verified"
	}

	// 3) Идемпотентность
	if k.isExecuted(msg.ID) {
		return qtypes.StatusExecuted, "replayed"
	}

	// 4) Rate-limit
	if exceeded, why := k.rateLimited(msg); exceeded {
		return qtypes.StatusRejected, "rate-limit: " + why
	}

	// 5) Выполнение по маршруту (пока заглушки)
	switch msg.Route {
	case qtypes.RouteTokenTransfer:
		// TODO: mint/unlock
	case qtypes.RouteContractCall:
		// TODO: IBC/contract invoke
	default:
		return qtypes.StatusRejected, "unsupported-route"
	}

	// 6) Отметить исполнение и событие
	k.markExecuted(msg.ID)
	k.emitEvent("bridge_execute", map[string]string{"msg_id": msg.ID})

	return qtypes.StatusExecuted, ""
}

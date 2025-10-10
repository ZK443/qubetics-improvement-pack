package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// Идемпотентное выполнение после успешного Verify.
// Никаких зависимостей от Cosmos SDK — чистая логика для CI.
func (k *Keeper) Execute(msg types.MsgExecute) (*types.MsgExecuteResponse, error) {
	// 1) Kill-switch
	if k.isPaused() {
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "paused"}, nil
	}

	// 2) Должен быть Verify
	if st := k.getStatusByID(msg.ID); st != types.StatusVerified {
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "not-verified"}, nil
	}

	// 3) Идемпотентность
	if k.isExecuted(msg.ID) {
		return &types.MsgExecuteResponse{Status: types.StatusExecuted, Reason: "replayed"}, nil
	}

	// 4) Rate-limit
	if exceeded, why := k.rateLimited(msg); exceeded {
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "rate-limit: " + why}, nil
	}

	// 5) Выполнение по маршруту (пока заглушки)
	switch msg.Route {
	case types.RouteTokenTransfer, types.RouteContractCall:
		// TODO: подключить реальный адаптер (mint/unlock/contract-call)
	default:
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "unsupported-route"}, nil
	}

	// 6) Отметить исполнение + событие
	k.markExecuted(msg.ID)
	k.emitEvent("bridge_execute", map[string]string{"msg_id": msg.ID})

	return &types.MsgExecuteResponse{Status: types.StatusExecuted}, nil
}

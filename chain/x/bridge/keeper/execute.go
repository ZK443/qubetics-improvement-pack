package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

func (k *Keeper) Execute(msg types.MsgExecute) (*types.MsgExecuteResponse, error) {
	// 1) Kill-switch
	if k.isPaused() {
		k.emitEvent(types.EventPausedBlock, map[string]string{"msg_id": msg.ID})
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "paused"}, nil
	}

	// 2) Должен быть Verify
	if st := k.getStatusByID(msg.ID); st != types.StatusVerified {
		k.emitEvent(types.EventExecuteDenied, map[string]string{"msg_id": msg.ID, "reason": "not-verified"})
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "not-verified"}, nil
	}

	// 3) Идемпотентность
	if k.isExecuted(msg.ID) {
		k.emitEvent(types.EventExecuteReplay, map[string]string{"msg_id": msg.ID})
		return &types.MsgExecuteResponse{Status: types.StatusExecuted, Reason: "replayed"}, nil
	}

	// 4) Rate-limit
	if exceeded, why := k.rateLimited(msg); exceeded {
		k.emitEvent(types.EventRateLimitHit, map[string]string{"msg_id": msg.ID, "why": why})
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "rate-limit: " + why}, nil
	}

	// 5) Выполнение по маршруту (пока заглушки)
	switch msg.Route {
	case types.RouteTokenTransfer, types.RouteContractCall:
	default:
		k.emitEvent(types.EventUnsupported, map[string]string{"msg_id": msg.ID, "route": string(msg.Route)})
		return &types.MsgExecuteResponse{Status: types.StatusRejected, Reason: "unsupported-route"}, nil
	}

	// 6) Отметить исполнение + событие
	k.markExecuted(msg.ID)
	k.emitEvent(types.EventExecuteOK, map[string]string{"msg_id": msg.ID})

	return &types.MsgExecuteResponse{Status: types.StatusExecuted}, nil
}

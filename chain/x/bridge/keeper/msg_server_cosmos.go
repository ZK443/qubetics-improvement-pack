//go:build cosmos

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	brtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// msgServer реализует bridge.v1.Msg (сгенерированный интерфейс совпадает по методам).
type msgServer struct {
	Keeper
}

func NewCosmosMsgServer(k Keeper) *msgServer { return &msgServer{Keeper: k} }

// VerifyProof: заглушка с установкой статуса Verified и событием.
func (s *msgServer) VerifyProof(ctx sdk.Context, msg *brtypes.MsgVerifyProof) (*brtypes.MsgVerifyProofResponse, error) {
	if s.isPaused(ctx) {
		s.emitEvent(ctx, "bridge.verify.denied", map[string]string{"proof_id": msg.Proof.ProofId, "why": "paused"})
		return &brtypes.MsgVerifyProofResponse{Status: "paused"}, nil
	}
	// TODO: валидация digest/подписи/доказательства
	s.SetStatus(ctx, msg.MessageId, brtypes.StatusVerified)
	s.emitEvent(ctx, "bridge.verify.ok", map[string]string{"msg_id": msg.MessageId})
	return &brtypes.MsgVerifyProofResponse{Status: "verified"}, nil
}

// Execute: маппинг из protobuf MsgExecute на внутренние типы и охранные проверки.
func (s *msgServer) Execute(ctx sdk.Context, msg *brtypes.MsgExecute) (*brtypes.MsgExecuteResponse, error) {
	if !s.CanExecute(ctx, msg.MessageId) {
		return &brtypes.MsgExecuteResponse{Status: "rejected"}, nil
	}

	// Маппинг маршрута
	route := brtypes.RouteTokenTransfer
	switch msg.Route {
	case brtypes.ExecRoute_EXEC_ROUTE_TOKEN_TRANSFER:
		route = brtypes.RouteTokenTransfer
	case brtypes.ExecRoute_EXEC_ROUTE_CONTRACT_CALL:
		route = brtypes.RouteContractCall
	default:
		s.emitEvent(ctx, brtypes.EventUnsupported, map[string]string{"msg_id": msg.MessageId})
		return &brtypes.MsgExecuteResponse{Status: "rejected"}, nil
	}

	// Rate-limit
	if limited, why := s.rateLimited(ctx, brtypes.MsgExecute{ // используем внутренний тип
		ID:    msg.MessageId,
		Route: route,
	}); limited {
		s.emitEvent(ctx, brtypes.EventRateLimitHit, map[string]string{"msg_id": msg.MessageId, "why": why})
		return &brtypes.MsgExecuteResponse{Status: "rate-limited"}, nil
	}

	// Mark executed
	s.MarkExecuted(ctx, msg.MessageId)
	s.emitEvent(ctx, brtypes.EventExecuteOK, map[string]string{"msg_id": msg.MessageId})
	return &brtypes.MsgExecuteResponse{Status: "executed"}, nil
}

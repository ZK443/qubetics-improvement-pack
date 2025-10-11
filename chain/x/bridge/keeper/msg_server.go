//go:build cosmos

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// MsgServer реализация поверх Cosmos Keeper.
type msgServer struct {
	Keeper
}

// NewMsgServer создаёт новый Cosmos MsgServer поверх Keeper.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// VerifyProof — проверка доказательства (заглушка для прототипа).
func (s *msgServer) VerifyProof(ctx sdk.Context, msg *types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
	if s.isPaused(ctx) {
		s.emitEvent(ctx, "bridge.verify.denied", map[string]string{"proof_id": msg.ProofId, "why": "paused"})
		return &types.MsgVerifyProofResponse{Status: "paused"}, nil
	}

	// (здесь будет реальная проверка подписи/доказательства)
	s.SetStatus(ctx, msg.MessageID, types.StatusVerified)
	s.emitEvent(ctx, "bridge.verify.ok", map[string]string{"msg_id": msg.MessageID})

	return &types.MsgVerifyProofResponse{Status: "verified"}, nil
}

// Execute — вызов после успешного Verify.
func (s *msgServer) Execute(ctx sdk.Context, msg *types.MsgExecute) (*types.MsgExecuteResponse, error) {
	if !s.CanExecute(ctx, msg.MessageId) {
		return &types.MsgExecuteResponse{Status: "rejected"}, nil
	}

	if limited, why := s.rateLimited(ctx, *msg); limited {
		s.emitEvent(ctx, types.EventRateLimitHit, map[string]string{"msg_id": msg.MessageId, "why": why})
		return &types.MsgExecuteResponse{Status: "rate-limited"}, nil
	}

	s.MarkExecuted(ctx, msg.MessageId)
	s.emitEvent(ctx, types.EventExecuteOK, map[string]string{"msg_id": msg.MessageId})
	return &types.MsgExecuteResponse{Status: "executed"}, nil
}

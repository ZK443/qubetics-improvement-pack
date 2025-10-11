//go:build !cosmos

package keeper

import "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// MsgServer (in-memory) — используется в CI и unit-тестах без SDK.
type MsgServer struct {
	k *Keeper
}

func NewMsgServer(k *Keeper) *MsgServer { return &MsgServer{k: k} }

func (s *MsgServer) VerifyProof(_ interface{}, msg *types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
	if s.k.isPaused() {
		s.k.emitEvent("bridge.verify.denied", map[string]string{"proof_id": msg.ProofId, "why": "paused"})
		return &types.MsgVerifyProofResponse{Status: "paused"}, nil
	}

	s.k.SetStatus(msg.MessageID, types.StatusVerified)
	s.k.emitEvent("bridge.verify.ok", map[string]string{"msg_id": msg.MessageID})
	return &types.MsgVerifyProofResponse{Status: "verified"}, nil
}

func (s *MsgServer) Execute(_ interface{}, msg *types.MsgExecute) (*types.MsgExecuteResponse, error) {
	out, _ := s.k.Execute(*msg)
	return out, nil
}

//go:build cosmos

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	brtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

type queryServer struct {
	k Keeper
}

func NewQueryServer(k Keeper) *queryServer { return &queryServer{k: k} }

// --- Query impl ---

func (s *queryServer) Params(ctx context.Context, _ *brtypes.ParamsRequest) (*brtypes.ParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	p := s.k.GetParams(sdkCtx)
	return &brtypes.ParamsResponse{
		GlobalPause:       p.GlobalPause,
		RateLimitAmount:   p.RateLimitAmount,
		RateLimitWindowMs: p.RateLimitWindowMs,
	}, nil
}

func (s *queryServer) Status(ctx context.Context, req *brtypes.StatusRequest) (*brtypes.StatusResponse, error) {
	if req == nil || req.MessageId == "" {
		return &brtypes.StatusResponse{Status: uint32(brtypes.StatusUnknown)}, nil
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	st := s.k.GetStatus(sdkCtx, req.MessageId)
	return &brtypes.StatusResponse{Status: uint32(st)}, nil
}

func (s *queryServer) Nonce(ctx context.Context, req *brtypes.NonceRequest) (*brtypes.NonceResponse, error) {
	if req == nil || req.Sender == "" {
		return &brtypes.NonceResponse{Nonce: 0}, nil
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	n := s.k.PeekNonce(sdkCtx, req.Sender)
	return &brtypes.NonceResponse{Nonce: n}, nil
}

func (s *queryServer) IsAllowed(ctx context.Context, req *brtypes.IsAllowedRequest) (*brtypes.IsAllowedResponse, error) {
	if req == nil || req.Addr == "" {
		return &brtypes.IsAllowedResponse{Allowed: false}, nil
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ok := s.k.IsAllowed(sdkCtx, req.Addr)
	return &brtypes.IsAllowedResponse{Allowed: ok}, nil
}

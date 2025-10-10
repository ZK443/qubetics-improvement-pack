package types

import "cosmossdk.io/errors"

var (
    ErrUnauthorized      = errors.Register(ModuleName, 1200, "unauthorized")
    ErrInvalidProof      = errors.Register(ModuleName, 1201, "invalid proof")
    ErrNotVerified       = errors.Register(ModuleName, 1202, "proof not verified")
    ErrProofAlreadyExists= errors.Register(ModuleName, 1203, "proof already exists")
    ErrAlreadyExecuted   = errors.Register(ModuleName, 1204, "already executed")
    ErrBridgePaused      = errors.Register(ModuleName, 1205, "bridge paused")
    ErrRateLimitExceeded = errors.Register(ModuleName, 1206, "rate limit exceeded")
)

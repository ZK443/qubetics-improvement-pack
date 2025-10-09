package types

import "cosmossdk.io/errors"

var (
    ErrProofAlreadyExists = errors.Register(ModuleName, 1100, "proof already exists")
    ErrAlreadyExecuted    = errors.Register(ModuleName, 1101, "message already executed")
)

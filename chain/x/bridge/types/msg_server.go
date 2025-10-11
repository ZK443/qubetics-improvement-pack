package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgServer — интерфейс обработчиков сообщений модуля bridge.
// Cosmos SDK автоматически связывает его с gRPC/Tx.
type MsgServer interface {
	VerifyProof(sdk.Context, *MsgVerifyProof) (*MsgVerifyProofResponse, error)
	Execute(sdk.Context, *MsgExecute) (*MsgExecuteResponse, error)
}

// Общие коды ошибок для MsgServer (единая карта отказов).
const (
	ErrUnauthorized   = "unauthorized"
	ErrNotVerified    = "not-verified"
	ErrAlreadyExecuted = "already-executed"
	ErrRateLimited    = "rate-limited"
	ErrPaused         = "paused"
	ErrUnsupported    = "unsupported-route"
	ErrInternal       = "internal"
)

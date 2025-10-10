package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgVerifyProof struct {
    Signer        sdk.AccAddress
    ProofId       string
    ProofData     []byte
    MessageID     string
    ExpectedDigest []byte
    Verifier      string
}

type MsgVerifyProofResponse struct {
    Status string
}

type MsgExecute struct {
    Executor  sdk.AccAddress
    MessageId string
    ProofId   string
    Amount    sdk.Int // для rate-limit примера
}

type MsgExecuteResponse struct {
    Status string
}

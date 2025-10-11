package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgVerifyProof = "verify_proof"
	TypeMsgExecute     = "execute"
	RouterKey          = ModuleName
)

// ---------- MsgVerifyProof ----------

type MsgVerifyProof struct {
	Sender         string `json:"sender" yaml:"sender"`             // bech32
	MsgID          string `json:"msg_id" yaml:"msg_id"`             // детермин. идентификатор
	ProofID        string `json:"proof_id" yaml:"proof_id"`         // внешний id доказательства (если есть)
	ProofData      []byte `json:"proof_data" yaml:"proof_data"`
	ExpectedDigest []byte `json:"expected_digest" yaml:"expected_digest"`
	Verifier       string `json:"verifier" yaml:"verifier"`         // произвольный идентификатор верификатора
	Meta           []byte `json:"meta" yaml:"meta"`                 // метаданные для binding
}

func (m *MsgVerifyProof) Route() string { return RouterKey }
func (m *MsgVerifyProof) Type() string  { return TypeMsgVerifyProof }
func (m *MsgVerifyProof) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}
func (m *MsgVerifyProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %v", err)
	}
	if len(m.MsgID) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty msg_id")
	}
	if len(m.ProofData) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty proof_data")
	}
	return nil
}
func (m *MsgVerifyProof) String() string { return fmt.Sprintf("VerifyProof{%s}", m.MsgID) }

type MsgVerifyProofResponse struct {
	Status string `json:"status" yaml:"status"`
}

// ---------- MsgExecute ----------

type MsgExecute struct {
	Sender  string `json:"sender" yaml:"sender"` // bech32
	MsgID   string `json:"msg_id" yaml:"msg_id"`
	ProofID string `json:"proof_id" yaml:"proof_id"`
	Route   Route  `json:"route" yaml:"route"`   // из types/store.go
	Meta    []byte `json:"meta" yaml:"meta"`
}

func (m *MsgExecute) Route() string { return RouterKey }
func (m *MsgExecute) Type() string  { return TypeMsgExecute }
func (m *MsgExecute) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{addr}
}
func (m *MsgExecute) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "sender: %v", err)
	}
	if len(m.MsgID) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty msg_id")
	}
	return nil
}

type MsgExecuteResponse struct {
	Status string `json:"status" yaml:"status"`
}

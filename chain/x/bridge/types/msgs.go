// SPDX-License-Identifier: MIT
package types

const (
	ActionSend    = "bridge_send"
	ActionVerify  = "bridge_verify"
	ActionExecute = "bridge_execute"
)

type Msg interface{ Type() string }

// MsgSend — заявка на межсетевой перенос/вызов.
type MsgSend struct {
	ID      string
	Nonce   uint64
	Source  ChainID
	Dest    ChainID
	Route   Route
	Payload []byte
	Sender  []byte
}
func (m MsgSend) Type() string { return ActionSend }

// MsgVerify — доказательство для ранее созданного сообщения.
type MsgVerify struct {
	MsgID string
	Proof Proof
}
func (m MsgVerify) Type() string { return ActionVerify }

// MsgExecute — попытка выполнить ранее верифицированное сообщение.
type MsgExecute struct {
	MsgID string
}
func (m MsgExecute) Type() string { return ActionExecute }

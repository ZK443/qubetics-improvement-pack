// SPDX-License-Identifier: MIT
// Package: chain/x/bridge/types
package types

// NB: это чистый каркас без Cosmos-SDK зависимостей.
// Он описывает формат сообщений и событий для моста.
// Реальные Msg* для Cosmos-SDK можно будет сгенерировать позднее.

const (
	// High-level actions
	ActionSend    = "bridge_send"
	ActionVerify  = "bridge_verify"
	ActionExecute = "bridge_execute"
)

// MsgSend — заявка на межсетевой перенос/вызов.
// Инициируется пользователем на исходной сети (или на Qubetics для исходящих).
type MsgSend struct {
	ID      string  // детерминированный hash заявки (может быть пустым — вычисляется на стороне клиента)
	Nonce   uint64  // monotonically increasing per (sender, route) — защита от повторов
	Source  ChainID // исходная цепь
	Dest    ChainID // целевая цепь
	Route   Route   // token-transfer / contract-call / etc
	Payload []byte  // полезная нагрузка (ABI/borsh/…)
	Sender  []byte  // адрес инициатора (raw-байты без привязки к конкретной сети)
}

// MsgVerify — доказательство для ранее созданного сообщения.
type MsgVerify struct {
	MsgID string // должен ссылаться на MsgSend.ID
	Proof Proof  // контейнер лёгкого/zk/SPV доказательства
}

// MsgExecute — попытка выполнить ранее верифицированное сообщение.
type MsgExecute struct {
	MsgID string // тот же идентификатор
	// Опционально можно добавить подписи исполнителей/релайеров, если политика этого требует.
}

// Helper-интерфейсы для будущей SDK-реализации.
// Они полезны для моков в тестах, пока нет полной интеграции.
type Msg interface {
	Type() string
}

func (m MsgSend) Type() string    { return ActionSend }
func (m MsgVerify) Type() string  { return ActionVerify }
func (m MsgExecute) Type() string { return ActionExecute }

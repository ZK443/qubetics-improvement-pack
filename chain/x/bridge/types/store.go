// SPDX-License-Identifier: MIT
// Package: chain/x/bridge/types
package types

// Каркас ключей хранилища и статусов сообщений моста.
// Минимум зависимостей, чтобы не мешать CI.

type Status uint8

const (
	StatusUnknown Status = iota
	StatusPending   // отправлено (Send), ждём Verify
	StatusVerified  // доказ-во принято (Verify ok)
	StatusExecuted  // выполнено (Execute ok, идемпотентно)
	StatusFailed    // ошибка выполнения
)

// Префиксы для KV.
var (
	KeyPrefixMsg    = []byte{0x10} // msg:ID -> MsgSend (raw)
	KeyPrefixStatus = []byte{0x11} // st:ID     -> Status
	KeyPrefixNonce  = []byte{0x12} // nc:<sender-route> -> uint64
	KeyPrefixProof  = []byte{0x13} // pf:ID     -> Proof blob
	KeyPrefixLimits = []byte{0x14} // rl:<route>-> RateLimitConfig
	KeyPrefixPause  = []byte{0x15} // ps:route  -> bool (paused)
)

// Хелперы ключей
func KeyMsg(id string) []byte       { return append(append([]byte{}, KeyPrefixMsg...), []byte(id)...) }
func KeyStatus(id string) []byte    { return append(append([]byte{}, KeyPrefixStatus...), []byte(id)...) }
func KeyNonce(route string) []byte  { return append(append([]byte{}, KeyPrefixNonce...), []byte(route)...) }
func KeyProof(id string) []byte     { return append(append([]byte{}, KeyPrefixProof...), []byte(id)...) }
func KeyLimit(route string) []byte  { return append(append([]byte{}, KeyPrefixLimits...), []byte(route)...) }
func KeyPause(route string) []byte  { return append(append([]byte{}, KeyPrefixPause...), []byte(route)...) }

// Базовые типы
type ChainID string
type Route string

const (
	RouteTokenTransfer Route = "token-transfer"
	RouteContractCall  Route = "contract-call"
)

type ProofKind string

const (
	ProofLight ProofKind = "light"
	ProofZK    ProofKind = "zk"
	ProofSPV   ProofKind = "spv"
)

type Proof struct {
	Kind  ProofKind
	Bytes []byte // формат определяется адаптером клиента (ETH/SOL/BTC)
}

type Message struct {
	ID     string
	Nonce  uint64
	Source ChainID
	Dest   ChainID
	Route  Route
	// payload/sender держим вне каркаса, см. MsgSend
}

type VerificationResult struct {
	Valid  bool
	Reason string
}

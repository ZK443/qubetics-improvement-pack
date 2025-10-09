// SPDX-License-Identifier: MIT
package client

import qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"

// Client — общий интерфейс верификаторов.
// Каждая сеть (ETH/SOL/BTC) реализует Verify() своим способом.
type Client interface {
	Network() string
	Verify(msg qtypes.Message, proof qtypes.Proof) qtypes.VerificationResult
}

// Registry — простейший реестр клиентов по имени сети.
type Registry struct{ m map[string]Client }

func NewRegistry() *Registry { return &Registry{m: make(map[string]Client)} }

func (r *Registry) Register(c Client) { r.m[c.Network()] = c }

func (r *Registry) Get(name string) (Client, bool) {
	c, ok := r.m[name]
	return c, ok
}

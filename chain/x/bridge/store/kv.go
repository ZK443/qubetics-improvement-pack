// SPDX-License-Identifier: MIT
package store

import "sync"

type KV struct {
	mu sync.RWMutex
	m  map[string][]byte
}

func New() *KV { return &KV{m: make(map[string][]byte)} }

func (s *KV) Get(k []byte) []byte {
	s.mu.RLock(); defer s.mu.RUnlock()
	v := s.m[string(k)]
	if v == nil { return nil }
	cp := make([]byte, len(v))
	copy(cp, v)
	return cp
}

func (s *KV) Set(k, v []byte) {
	s.mu.Lock(); defer s.mu.Unlock()
	cp := make([]byte, len(v))
	copy(cp, v)
	s.m[string(k)] = cp
}

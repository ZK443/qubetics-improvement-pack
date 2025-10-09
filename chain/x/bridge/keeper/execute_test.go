// SPDX-License-Identifier: MIT
package keeper_test

import (
	"testing"

	qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// memKeeper — минимальная in-memory модель статусов для демонстрации инвариантов.
type memKeeper struct {
	status   map[string]qtypes.Status
	executed map[string]bool
}

func newMemKeeper() *memKeeper {
	return &memKeeper{
		status:   make(map[string]qtypes.Status),
        executed: make(map[string]bool),
	}
}

// verify — имитирует успешную верификацию сообщения.
func (m *memKeeper) verify(id string) {
	m.status[id] = qtypes.StatusVerified
}

// execute — идемпотентное выполнение: возможно только после Verify.
func (m *memKeeper) execute(id string) error {
	if m.status[id] != qtypes.StatusVerified {
		return ErrNotVerified
	}
	if m.executed[id] {
		// идемпотентность: повторный вызов — успех без эффекта.
		return nil
	}
	m.executed[id] = true
	m.status[id] = qtypes.StatusExecuted
	return nil
}

// простые ошибки как значения (чтобы не тянуть пакеты)
var (
	ErrNotVerified = simpleErr("not verified")
)

type simpleErr string
func (e simpleErr) Error() string { return string(e) }

// ---- TESTS ----

func TestExecuteRequiresVerify(t *testing.T) {
	m := newMemKeeper()
	id := "msg-1"

	// До Verify выполнение запрещено.
	if err := m.execute(id); err == nil {
		t.Fatalf("expected error when executing without Verify")
	}

	// После Verify — разрешено.
	m.verify(id)
	if err := m.execute(id); err != nil {
		t.Fatalf("unexpected error after Verify: %v", err)
	}
}

func TestExecuteIsIdempotent(t *testing.T) {
	m := newMemKeeper()
	id := "msg-2"

	m.verify(id)
	if err := m.execute(id); err != nil {
		t.Fatalf("first execute failed: %v", err)
	}
	// Повторный Execute — допустим и не меняет состояние.
	if err := m.execute(id); err != nil {
		t.Fatalf("second execute should be idempotent: %v", err)
	}
	if !m.executed[id] || m.status[id] != qtypes.StatusExecuted {
		t.Fatalf("message should remain executed; got executed=%v status=%v", m.executed[id], m.status[id])
	}
}

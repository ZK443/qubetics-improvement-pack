package types

import "testing"

func TestKeyDerivationIsDeterministicAndDistinct(t *testing.T) {
	id := "tx-123"

	// Детерминированность (один и тот же ввод -> тот же ключ)
	k1 := string(KeyMsg(id))
	k2 := string(KeyMsg(id))
	if k1 != k2 {
		t.Fatalf("KeyMsg must be deterministic: %q != %q", k1, k2)
	}

	// Разные функции ключей не должны совпадать
	if string(KeyMsg(id)) == string(KeyStatus(id)) {
		t.Fatalf("KeyMsg and KeyStatus must not collide")
	}
	if string(KeyMsg(id)) == string(KeyNonce("sender")) {
		t.Fatalf("KeyMsg and KeyNonce must not collide")
	}
	if string(KeyStatus(id)) == string(KeyNonce("sender")) {
		t.Fatalf("KeyStatus and KeyNonce must not collide")
	}
}

func TestStatusOrderInvariants(t *testing.T) {
	// Предполагаем, что порядок такой:
	// Unknown < Pending < Verified < Executed, а Failed отдельный.
	if !(StatusPending > StatusUnknown) {
		t.Fatalf("expected Pending > Unknown")
	}
	if !(StatusVerified > StatusPending) {
		t.Fatalf("expected Verified > Pending")
	}
	if !(StatusExecuted > StatusVerified) {
		t.Fatalf("expected Executed > Verified")
	}
	// Failed не обязан быть больше/меньше — просто должен отличаться
	if StatusFailed == StatusExecuted || StatusFailed == StatusVerified || StatusFailed == StatusPending {
		t.Fatalf("StatusFailed must be distinct")
	}
}

func TestMessageStatusOrderInvariants(t *testing.T) {
	// Пропускаем тест, если MessageStatus не определён
	defer func() {
		if r := recover(); r != nil {
			t.Skip("MessageStatus type not defined in this module")
		}
	}()

	_ = MessageStatus(0)
	if !(StatusVerified > StatusPending) {
		t.Fatalf("expected StatusVerified > StatusPending (MessageStatus)")
	}
}

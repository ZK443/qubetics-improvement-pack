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

// Инварианты порядка статусов (минимальный набор для прототипа).
func TestStatusOrderInvariants(t *testing.T) {
	if !(StatusUnknown < StatusPending &&
		StatusPending < StatusVerified &&
		StatusVerified < StatusExecuted &&
		StatusExecuted < StatusFailed) {
		t.Fatalf("status order invariants failed")
	}
}

package types

import "testing"

func TestKeyDerivationIsDeterministicAndDistinct(t *testing.T) {
	id := "tx-123"

	// детерминированность
	k1 := string(KeyMsg(id))
	k2 := string(KeyMsg(id))
	if k1 != k2 {
		t.Fatalf("KeyMsg must be deterministic: %q != %q", k1, k2)
	}

	// коллизий между разными ключами быть не должно
	if string(KeyMsg(id)) == string(KeyStatus(id)) {
		t.Fatalf("KeyMsg and KeyStatus must not collide")
	}
	route := RouteTokenTransfer // <— корректный аргумент для KeyNonce
	if string(KeyMsg(id)) == string(KeyNonce(route)) {
		t.Fatalf("KeyMsg and KeyNonce must not collide")
	}
	if string(KeyStatus(id)) == string(KeyNonce(route)) {
		t.Fatalf("KeyStatus and KeyNonce must not collide")
	}
}

// Инварианты для единого Status
func TestStatusOrderInvariants(t *testing.T) {
	if !(StatusUnknown < StatusPending &&
		StatusPending < StatusVerified &&
		StatusVerified < StatusExecuted &&
		StatusExecuted < StatusFailed) {
		t.Fatalf("status order invariants failed")
	}
}

/*
НЕ добавляем TestMessageStatusOrderInvariants без отдельного типа MessageStatus:
он упадёт на этапе компиляции, если типа нет.

Если/когда появится отдельный enum, добавим отдельный файл с билд-тегом:

//go:build msgstatus

package types

import "testing"

func TestMessageStatusOrderInvariants(t *testing.T) {
	if !(MessageStatusVerified > MessageStatusPending) {
		t.Fatalf("expected MessageStatusVerified > MessageStatusPending")
	}
}
*/

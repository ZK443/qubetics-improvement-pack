package keeper

import (
	"testing"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func TestStatusCRUDAndCanExecute(t *testing.T) {
	k := NewKeeper()
	id := "msg-1"

	if got := k.GetStatus(id); got != types.StatusUnknown {
		t.Fatalf("default status should be unknown")
	}
	k.SetStatus(id, types.StatusVerified)

	if !k.CanExecute(id) {
		t.Fatalf("should be executable after verified and no pause")
	}

	k.MarkExecuted(id)
	if k.CanExecute(id) {
		t.Fatalf("must not allow double execute")
	}
	if got := k.GetStatus(id); got != types.StatusExecuted {
		t.Fatalf("status should be executed")
	}
}

func TestNonceMonotonic(t *testing.T) {
	k := NewKeeper()
	addr := "sender1"
	if k.PeekNonce(addr) != 0 {
		t.Fatalf("default nonce must be 0")
	}
	n1 := k.NextNonce(addr)
	n2 := k.NextNonce(addr)
	if !(n1 == 1 && n2 == 2 && k.PeekNonce(addr) == 2) {
		t.Fatalf("nonce must be monotonic: got %d then %d, peek %d", n1, n2, k.PeekNonce(addr))
	}
}

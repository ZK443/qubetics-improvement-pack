package keeper

import (
	"testing"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func TestKeeper_ParamsCRUD(t *testing.T) {
	k := NewKeeper()

	// defaults
	def := k.GetParams()
	if def.RateLimitAmount == 0 {
		t.Fatalf("default RateLimitAmount must be > 0")
	}

	// set/get
	newP := def
	newP.GlobalPause = true
	if err := k.SetParams(newP); err != nil {
		t.Fatalf("unexpected SetParams error: %v", err)
	}
	got := k.GetParams()
	if !got.GlobalPause {
		t.Fatalf("expected GlobalPause=true")
	}
}

func TestKeeper_ACL(t *testing.T) {
	k := NewKeeper()
	addr := "relayer1"
	if k.IsAllowed(addr) {
		t.Fatalf("not allowed by default")
	}
	k.SetAllowed(addr, true)
	if !k.IsAllowed(addr) {
		t.Fatalf("expected allowed after SetAllowed")
	}
}

package keeper

import (
	"testing"

	qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func TestMsgServer_ExecuteAndVerify(t *testing.T) {
	k := NewKeeper()
	srv := NewMsgServer(k)

	// VerifyProof
	resp1, _ := srv.VerifyProof(nil, &qtypes.MsgVerifyProof{MessageID: "m1"})
	if resp1.Status != "verified" {
		t.Fatalf("expected verified, got %v", resp1.Status)
	}

	// Execute
	msg := &qtypes.MsgExecute{MessageId: "m1", Route: qtypes.RouteTokenTransfer}
	resp2, _ := srv.Execute(nil, msg)
	if resp2.Status != "executed" {
		t.Fatalf("expected executed, got %v", resp2.Status)
	}

	// Проверка событий
	evs := k.Events()
	if len(evs) < 2 {
		t.Fatalf("expected at least 2 events, got %d", len(evs))
	}
}

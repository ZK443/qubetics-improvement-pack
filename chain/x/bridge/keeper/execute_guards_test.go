package keeper

import (
	"testing"

	qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func newKeeperForTests() *Keeper {
	k := NewKeeper()
	// более агрессивные лимиты для тестов
	p := k.GetParams()
	p.RateLimitWindowMs = 1_000 // 1s окно
	p.RateLimitAmount = 2       // максимум 2 выполнения за окно
	_ = k.SetParams(p)
	return k
}

func TestExecute_NotVerified(t *testing.T) {
	k := newKeeperForTests()
	msg := qtypes.MsgExecute{ID: "m1", Route: qtypes.RouteTokenTransfer}

	resp, err := k.Execute(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != qtypes.StatusRejected || resp.Reason != "not-verified" {
		t.Fatalf("expected not-verified rejection, got %+v", resp)
	}
}

func TestExecute_Idempotent(t *testing.T) {
	k := newKeeperForTests()
	id := "m2"
	k.SetStatus(id, qtypes.StatusVerified)

	msg := qtypes.MsgExecute{ID: id, Route: qtypes.RouteTokenTransfer}

	// Первое исполнение — OK
	resp, _ := k.Execute(msg)
	if resp.Status != qtypes.StatusExecuted {
		t.Fatalf("first execute should succeed, got %+v", resp)
	}

	// Повтор — executed/replayed
	resp, _ = k.Execute(msg)
	if resp.Status != qtypes.StatusExecuted || resp.Reason != "replayed" {
		t.Fatalf("second execute should be idempotent, got %+v", resp)
	}
}

func TestExecute_Paused(t *testing.T) {
	k := newKeeperForTests()
	// Включить паузу параметром
	p := k.GetParams()
	p.GlobalPause = true
	_ = k.SetParams(p)

	msg := qtypes.MsgExecute{ID: "m3", Route: qtypes.RouteTokenTransfer}
	resp, _ := k.Execute(msg)
	if resp.Status != qtypes.StatusRejected || resp.Reason != "paused" {
		t.Fatalf("expected paused rejection, got %+v", resp)
	}
}

func TestExecute_RateLimitWindow(t *testing.T) {
	k := newKeeperForTests()

	// Подменить «часы», чтобы контролировать окно
	cur := int64(1_000_000)
	k.now = func() int64 { return cur }

	route := qtypes.RouteTokenTransfer

	// Подготовить три разных сообщения и все пометить Verified
	for _, id := range []string{"r1", "r2", "r3"} {
		k.SetStatus(id, qtypes.StatusVerified)
	}

	// Первые два — влезают в окно
	for _, id := range []string{"r1", "r2"} {
		resp, _ := k.Execute(qtypes.MsgExecute{ID: id, Route: route})
		if resp.Status != qtypes.StatusExecuted {
			t.Fatalf("expected executed for %s, got %+v", id, resp)
		}
	}

	// Третий — должен упасть по rate-limit
	resp, _ := k.Execute(qtypes.MsgExecute{ID: "r3", Route: route})
	if resp.Status != qtypes.StatusRejected || resp.Reason[:10] != "rate-limit" {
		t.Fatalf("expected rate-limit rejection, got %+v", resp)
	}

	// Сместить время за пределы окна и убедиться, что снова можно
	cur += 1_200 // > 1s окна
	k.SetStatus("r4", qtypes.StatusVerified)
	resp, _ = k.Execute(qtypes.MsgExecute{ID: "r4", Route: route})
	if resp.Status != qtypes.StatusExecuted {
		t.Fatalf("expected executed after window reset, got %+v", resp)
	}
}

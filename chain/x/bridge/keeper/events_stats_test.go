package keeper

import (
	"testing"

	qtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func TestEvents_And_Stats(t *testing.T) {
	k := NewKeeper()

	// Настроить лимиты, чтобы сработал rate-limit.
	p := k.GetParams()
	p.RateLimitWindowMs = 1_000
	p.RateLimitAmount = 1
	_ = k.SetParams(p)

	// Подготовка: две попытки по одному ID + одна лишняя по лимиту.
	k.SetStatus("v1", qtypes.StatusVerified)

	// 1) Успешное исполнение
	_, _ = k.Execute(qtypes.MsgExecute{ID: "v1", Route: qtypes.RouteTokenTransfer})

	// 2) Повтор — идемпотентность
	_, _ = k.Execute(qtypes.MsgExecute{ID: "v1", Route: qtypes.RouteTokenTransfer})

	// 3) Вторая транзакция в том же окне — rate-limit
	k.SetStatus("v2", qtypes.StatusVerified)
	_, _ = k.Execute(qtypes.MsgExecute{ID: "v2", Route: qtypes.RouteTokenTransfer})

	evs := k.Events()
	if len(evs) == 0 {
		t.Fatalf("expected events to be recorded")
	}

	var ok, replay, rl uint64
	for _, e := range evs {
		switch e.Name {
		case qtypes.EventExecuteOK:
			ok++
		case qtypes.EventExecuteReplay:
			replay++
		case qtypes.EventRateLimitHit:
			rl++
		}
	}

	if ok != 1 || replay != 1 || rl != 1 {
		t.Fatalf("unexpected counters from events: ok=%d replay=%d rl=%d", ok, replay, rl)
	}

	stats := k.GetStats()
	if !(stats.Executed == 1 && stats.Replayed == 1 && stats.RateLimit == 1) {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

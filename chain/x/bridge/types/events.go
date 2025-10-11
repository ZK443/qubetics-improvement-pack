package types

// Имя события и произвольные атрибуты (строковые для упрощения).
type Event struct {
	Name  string
	Attrs map[string]string
}

// Имена событий модуля bridge.
// Держите их стабильными (используются в тестах/дашбордах).
const (
	EventExecuteOK      = "bridge.execute.ok"
	EventExecuteDenied  = "bridge.execute.denied"
	EventExecuteReplay  = "bridge.execute.replay"
	EventRateLimitHit   = "bridge.ratelimit.hit"
	EventPausedBlock    = "bridge.paused"
	EventUnsupported    = "bridge.route.unsupported"
)

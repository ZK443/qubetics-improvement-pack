package store

// Префиксы для счётчиков rate-limit по маршруту.
// bridge/rl/count/<route>  -> uint64 (число исполнений в окне)
// bridge/rl/until/<route>  -> int64  (unix ms, конец окна)
var (
	KeyRLCountPrefix = []byte("bridge/rl/count/")
	KeyRLUntilPrefix = []byte("bridge/rl/until/")
)

func KeyRLCount(route string) []byte { return append(append([]byte{}, KeyRLCountPrefix...), []byte(route)...) }
func KeyRLUntil(route string) []byte { return append(append([]byte{}, KeyRLUntilPrefix...), []byte(route)...) }

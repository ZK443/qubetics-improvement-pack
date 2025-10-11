package types

// Статус обработки сообщения моста.
type Status uint8

const (
	StatusUnknown Status = iota // создано/неизвестно
	StatusPending               // отправлено (Send), ждём Verify
	StatusVerified              // доказ-во принято (Verify ok)
	StatusExecuted              // выполнено (Execute ok)
	StatusRejected              // отклонено (rate-limit/вето/ошибка)
	StatusFailed                // внутренняя ошибка выполнения
)

// Для совместимости с тестами: MessageStatus == Status.
type MessageStatus = Status

// Маршруты выполнения (что делает Execute).
type Route string

const (
	RouteTokenTransfer Route = "token-transfer"
	RouteContractCall  Route = "contract-call"
)

// ---- KV-ключи (префиксы) ----
// Простые строки для наглядности; Cosmos-префиксы можно добавить позже.
var (
	KeyPrefixMsg     = []byte("bridge/msg/")      // msgID -> raw blob
	KeyPrefixStatus  = []byte("bridge/status/")   // msgID -> MessageStatus
	KeyPrefixNonce   = []byte("bridge/nonce/")    // sender -> uint64
	KeyPrefixRateLim = []byte("bridge/rl/route/") // route -> RateLimitConfig
	KeyPrefixPause   = []byte("bridge/kill")      // глобальная пауза (bool)
)

// Хелперы по формированию ключей.
func KeyMsg(id string) []byte        { return append([]byte{}, append(KeyPrefixMsg, []byte(id)...)...) }
func KeyStatus(id string) []byte     { return append([]byte{}, append(KeyPrefixStatus, []byte(id)...)...) }
func KeyNonce(sender string) []byte  { return append([]byte{}, append(KeyPrefixNonce, []byte(sender)...)...) }
func KeyRateLimit(route Route) []byte { return append([]byte{}, append(KeyPrefixRateLim, []byte(route)...)...) }
func KeyPause() []byte               { return append([]byte{}, KeyPrefixPause...) }

package types

// Статус обработки сообщения моста.
type Status uint8

const (
	StatusUnknown Status = iota // создано/неизвестно
	StatusPending               // отправлено (Send), ждём Verify
	StatusVerified              // доказ-во принято (Verify ok)
	StatusExecuted              // выполнено (Execute ok)
	StatusRejected              // отклонено (rate-limit/вито/ошибка)
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

// Сообщение на выполнение.
type MsgExecute struct {
	ID     string
	Route  Route
	Sender string
}

// Результат выполнения.
type MsgExecuteResponse struct {
	Status Status
	Reason string // текстовая причина отклонения (если есть)
}

// ---- KV-ключи (префиксы) ----
// Преднамеренно простые, чтобы не тащить Cosmos SDK в прототип.
var (
	KeyPrefixMsg      = []byte("bridge/msg/")      // msgID -> raw blob
	KeyPrefixStatus   = []byte("bridge/status/")   // msgID -> MessageStatus
	KeyPrefixNonce    = []byte("bridge/nonce/")    // sender -> uint64
	KeyPrefixRateLim  = []byte("bridge/rl/route/") // route -> RateLimitConfig
	KeyPrefixPause    = []byte("bridge/kill")      // глобальная пауза (bool)
)

// Хелперы по формированию ключей.
func KeyMsg(id string) []byte     { return append([]byte{}, append(KeyPrefixMsg, []byte(id)...)...) }
func KeyStatus(id string) []byte  { return append([]byte{}, append(KeyPrefixStatus, []byte(id)...)...) }
func KeyNonce(sender string) []byte {
	return append([]byte{}, append(KeyPrefixNonce, []byte(sender)...)...)
}
func KeyRateLimit(route Route) []byte {
	return append([]byte{}, append(KeyPrefixRateLim, []byte(route)...)...)
}
func KeyPause() []byte { return append([]byte{}, KeyPrefixPause...) }

package types

// Параметры модуля моста.
// Без SDK-зависимостей: только примитивные типы.
//
// Примечание про время окна: duration в мс, чтобы не тянуть time.Duration.
type Params struct {
	// Глобальная пауза исполнения (kill-switch).
	GlobalPause bool `json:"global_pause" yaml:"global_pause"`

	// Rate-limit для входящих Execute по маршруту (общий лимит, без ACL-детализации).
	RateLimitWindowMs uint64 `json:"rate_limit_window_ms" yaml:"rate_limit_window_ms"` // окно в мс
	RateLimitAmount   uint64 `json:"rate_limit_amount"   yaml:"rate_limit_amount"`     // допустимое количество за окно

	// Простейший ACL: карта "адрес (строка) → разрешён/запрещён".
	// Слой Keeper решит, что считать ключом (bech32 и т.п.).
	ACL map[string]bool `json:"acl" yaml:"acl"`
}

// Значения по умолчанию (безопасные, консервативные).
func DefaultParams() Params {
	return Params{
		GlobalPause:       false,
		RateLimitWindowMs: 10_000, // 10 секунд
		RateLimitAmount:   100,    // на 10 секунд
		ACL:               map[string]bool{}, // пусто = «никто явно не разрешён»
	}
}

// Базовая валидация значений.
func (p Params) Validate() error {
	// amount==0 разрешает всё; для прототипа считаем это ошибкой конфигурации
	// чтобы не отключить защиту случайно.
	if p.RateLimitAmount == 0 {
		return ErrInvalidRateLimitAmount
	}
	// Окно меньше 100 мс — малополезно и может быть источником флапов.
	if p.RateLimitWindowMs < 100 {
		return ErrInvalidRateLimitWindow
	}
	return nil
}

// Локальные ошибки (без SDK).
var (
	ErrInvalidRateLimitAmount = simpleError("invalid rate-limit amount (must be > 0)")
	ErrInvalidRateLimitWindow = simpleError("invalid rate-limit window (must be >= 100ms)")
)

// простейший тип ошибки, чтобы не тянуть fmt/errors.
// (можно заменить на errors.New, если хотите)
type simpleError string
func (e simpleError) Error() string { return string(e) }

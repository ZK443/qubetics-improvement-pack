// SPDX-License-Identifier: MIT
package keeper

// RateLimitConfig — простая модель лимитов на маршрут.
type RateLimitConfig struct {
	WindowSeconds uint32
	MaxUnits      uint64 // единицы зависят от маршрута (сумма/шт.)
}

func (k Keeper) isPaused(route string) bool {
	// TODO: читать флаг kill-switch из стора
	return false
}

func (k Keeper) rateLimited(route string) (bool, string) {
	// TODO: проверка окна/квоты (например, per token per day)
	return false, ""
}

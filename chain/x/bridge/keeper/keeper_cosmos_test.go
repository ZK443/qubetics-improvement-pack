//go:build cosmos

package keeper

import "testing"

// Интеграционные тесты Cosmos-keeper требуют стора/контекста SDK.
// Заглушка, чтобы зафиксировать намерение и не ломать обычный CI.
func TestCosmosKeeper_RateLimitAndEvents(t *testing.T) {
	t.Skip("requires Cosmos SDK context/store scaffolding; to be implemented with app setup")
}

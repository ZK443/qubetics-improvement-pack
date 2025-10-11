package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// InitGenesis — минимальная заглушка для skeleton-приложения.
// На S6.T2 сюда будет добавлена загрузка Params/ACL/статусов из genesis.
func InitGenesis(_ sdk.Context) {
	_ = types.DefaultParams // якорь на пакет types, чтобы сохранить зависимость
}

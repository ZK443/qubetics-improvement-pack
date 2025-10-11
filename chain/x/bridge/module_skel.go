//go:build cosmos

package bridge

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
)

var _ module.AppModule = (*AppModule)(nil)
var _ module.HasServices = (*AppModule)(nil)

type AppModule struct {
	k keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule { return AppModule{k: k} }

// RegisterServices — заглушка: регистрация gRPC-сервисов будет добавлена на S6.T2.
func (am AppModule) RegisterServices(conf module.Configurator) { _ = conf }

// интерфейсы модуля, не используемые в skeleton-приложении
func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}
func (AppModule) Name() string        { return "bridge" }
func (AppModule) InitGenesis(context.Context, module.JSONCodec, []byte) {}
func (AppModule) ExportGenesis(context.Context, module.JSONCodec) []byte { return nil }

//go:build cosmos

package bridge

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// ---------- AppModuleBasic ----------

type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

func (AppModuleBasic) Name() string { return types.ModuleName }

// no-op: держим сигнатуру для совместимости
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (AppModuleBasic) DefaultGenesis(cdc module.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc module.JSONCodec, _ module.ClientTxContext, bz json.RawMessage) error {
	var gs types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return err
	}
	return gs.Validate()
}

// ---------- AppModule ----------

type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

var _ module.AppModule = AppModule{}
var _ module.HasGenesis = AppModule{}
var _ module.HasServices = AppModule{}

func NewAppModule(k keeper.Keeper) AppModule { return AppModule{k: k} }

// gRPC/MsgServer будет подключён на S7 (здесь — заглушка).
func (am AppModule) RegisterServices(conf module.Configurator) { _ = conf }

// InitGenesis — загрузка параметров и ACL через Keeper.
func (am AppModule) InitGenesis(ctx context.Context, cdc module.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if len(data) == 0 {
		gs = *types.DefaultGenesis()
	} else {
		// игнорируем ошибку: на этапе ValidateGenesis она уже проверялась
		_ = cdc.UnmarshalJSON(data, &gs)
	}
	if err := gs.Validate(); err != nil {
		// при ошибке — не падаем в skeleton-режиме
		return
	}

	sCtx := sdk.UnwrapSDKContext(ctx)

	// Параметры
	_ = am.k.SetParams(sCtx, gs.Params)

	// ACL
	for addr, allowed := range gs.ACL {
		am.k.SetAllowed(sCtx, addr, allowed)
	}
}

// ExportGenesis — обратное преобразование из Keeper в GenesisState.
func (am AppModule) ExportGenesis(ctx context.Context, cdc module.JSONCodec) json.RawMessage {
	sCtx := sdk.UnwrapSDKContext(ctx)
	ps := am.k.GetParams(sCtx)

	gs := types.GenesisState{
		GlobalPause: ps.GlobalPause,
		Params:      ps,
		ACL:         map[string]bool{}, // при необходимости выгрузки полного ACL — добавить обход KV
	}
	return cdc.MustMarshalJSON(&gs)
}

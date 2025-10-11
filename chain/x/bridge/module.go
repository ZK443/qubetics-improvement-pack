//go:build cosmos

package bridge

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
	"github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// ---------- AppModuleBasic ----------

type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {} // необязательно

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

// RegisterServices — пока пусто (gRPC MsgServer привяжется на S7 вместе с proto).
func (am AppModule) RegisterServices(conf module.Configurator) { _ = conf }

// Genesis: загрузка/экспорт через Keeper.
func (am AppModule) InitGenesis(ctx context.Context, cdc module.JSONCodec, data json.RawMessage) {
	var gs types.GenesisState
	if len(data) == 0 {
		gs = *types.DefaultGenesis()
	} else {
		_ = cdc.UnmarshalJSON(data, &gs) // безопасно: ошибки уже валидировались на App init
	}
	// синхронизация и валидация
	if err := gs.Validate(); err == nil {
		// применить параметры
		_ = am.k.SetParams(sdk.UnwrapSDKContext(ctx), gs.Params)
		// применить ACL
		sCtx := sdk.UnwrapSDKContext(ctx)
		for addr, allowed := range gs.ACL {
			am.k.SetAllowed(sCtx, addr, allowed)
		}
	}
}

func (am AppModule) ExportGenesis(ctx context.Context, cdc module.JSONCodec) json.RawMessage {
	sCtx := sdk.UnwrapSDKContext(ctx)
	ps := am.k.GetParams(sCtx)
	// восстановить флаг pause из Params
	gs := types.GenesisState{
		GlobalPause: ps.GlobalPause,
		Params:      ps,
		ACL:         map[string]bool{}, // на S7 можно выгрузить полный ACL, если потребуется
	}
	return cdc.MustMarshalJSON(&gs)
}

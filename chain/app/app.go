package app

import (
	"errors"

	cosmolog "cosmossdk.io/log"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	bridgemodule "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge"
	bridgekeeper "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
	bridgetypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// QubeticsApp — минимальное Cosmos-приложение с интегрированным модулем bridge.
type QubeticsApp struct {
	*baseapp.BaseApp
	BridgeKeeper bridgekeeper.Keeper
	keys         map[string]*storetypes.KVStoreKey
	mm           *module.Manager
}

func NewQubeticsApp() *QubeticsApp {
	logger := cosmolog.NewNopLogger()
	db := dbm.NewMemDB()

	bApp := baseapp.NewBaseApp("qubetics", logger, db, noopTxDecoder)

	// stores
	bridgeKey := storetypes.NewKVStoreKey(bridgetypes.ModuleName)
	kvKeys := map[string]*storetypes.KVStoreKey{
		bridgetypes.ModuleName: bridgeKey,
	}

	bApp.MountKVStores(kvKeys)
	if err := bApp.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// keepers
	bridgeKeeper := bridgekeeper.NewKeeper(nil, bridgeKey, paramtypes.Subspace{})

	// module manager
	mm := module.NewManager(
		bridgemodule.NewAppModule(bridgeKeeper),
	)

	// order
	mm.SetOrderInitGenesis(bridgetypes.ModuleName)
	mm.SetOrderBeginBlockers(bridgetypes.ModuleName)
	mm.SetOrderEndBlockers(bridgetypes.ModuleName)

	// register gRPC services (MsgServer/Query) via Configurator
	conf := module.NewConfigurator(appCodec(), bApp.MsgServiceRouter(), bApp.GRPCQueryRouter())
	mm.RegisterServices(conf)

	app := &QubeticsApp{
		BaseApp:      bApp,
		keys:         kvKeys,
		BridgeKeeper: bridgeKeeper,
		mm:           mm,
	}

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	return app
}

// --- стандартные ABCI хуки ---

func (app *QubeticsApp) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.mm.BeginBlock(ctx)
	return abci.ResponseBeginBlock{}
}

func (app *QubeticsApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	app.mm.EndBlock(ctx)
	return abci.ResponseEndBlock{}
}

func (app *QubeticsApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	app.mm.InitGenesis(ctx, app.appCodec(), req.AppStateBytes)
	return abci.ResponseInitChain{}
}

// minimal codec (interface registry without registrations is sufficient for genesis JSON here)
func (app *QubeticsApp) appCodec() codec.Codec {
	return appCodec()
}

func appCodec() codec.Codec {
	ir := codectypes.NewInterfaceRegistry()
	return codec.NewProtoCodec(ir)
}

func noopTxDecoder(_ []byte) (sdk.Tx, error) {
	return nil, errors.New("tx decoding not configured for QubeticsApp skeleton")
}

type noopBankKeeper struct{}

func (noopBankKeeper) SendCoinsFromModuleToAccount(
	_ sdk.Context, _ string, _ sdk.AccAddress, _ sdk.Coins,
) error { return nil }

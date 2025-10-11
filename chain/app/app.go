package app

import (
	"errors"

	cosmolog "cosmossdk.io/log"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

	// --- хранилище ---
	bridgeKey := storetypes.NewKVStoreKey(bridgetypes.ModuleName)
	kvKeys := map[string]*storetypes.KVStoreKey{
		bridgetypes.ModuleName: bridgeKey,
	}

	bApp.MountKVStores(kvKeys)
	if err := bApp.LoadLatestVersion(); err != nil {
		panic(err)
	}

	// --- Keeper модуля bridge ---
	bridgeKeeper := bridgekeeper.NewKeeper(bridgeKey, noopBankKeeper{})

	// --- Менеджер модулей ---
	cdc := codec.NewLegacyAmino()
	mm := module.NewManager(
		bridgemodule.NewAppModule(bridgeKeeper),
	)

	// --- Привязка модулей к жизненному циклу приложения ---
	mm.SetOrderInitGenesis(bridgetypes.ModuleName)
	mm.SetOrderBeginBlockers(bridgetypes.ModuleName)
	mm.SetOrderEndBlockers(bridgetypes.ModuleName)

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

// --- утилиты ---

func (app *QubeticsApp) appCodec() codec.Codec {
	return codec.NewProtoCodec(nil)
}

func noopTxDecoder(_ []byte) (sdk.Tx, error) {
	return nil, errors.New("tx decoding not configured for QubeticsApp skeleton")
}

type noopBankKeeper struct{}

func (noopBankKeeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	return nil
}

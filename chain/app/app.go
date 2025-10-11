package app

import (
	"errors"

	cosmolog "cosmossdk.io/log"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"


	bridgemodule "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge"
	bridgekeeper "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
	bridgetypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

type QubeticsApp struct {
	*baseapp.BaseApp
	BridgeKeeper bridgekeeper.Keeper
	keys         map[string]*storetypes.KVStoreKey
}

func NewQubeticsApp() *QubeticsApp {
	logger := cosmolog.NewNopLogger()
	db := dbm.NewMemDB()

	bApp := baseapp.NewBaseApp("qubetics", logger, db, noopTxDecoder)

	// Инициализация хранилища для bridge-модуля
	bridgeKey := storetypes.NewKVStoreKey(bridgetypes.ModuleName)
	kvKeys := map[string]*storetypes.KVStoreKey{
		bridgetypes.ModuleName: bridgeKey,
	}

	bApp.MountKVStores(kvKeys)
	if err := bApp.LoadLatestVersion(); err != nil {
		panic(err)
	}
	
	app := &QubeticsApp{
		BaseApp:      bApp,
		keys:         kvKeys,
		BridgeKeeper: bridgekeeper.NewKeeper(nil, bridgeKey, paramtypes.Subspace{}),
	}

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	
	return app
}

func (app *QubeticsApp) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (app *QubeticsApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (app *QubeticsApp) InitChainer(ctx sdk.Context, _ abci.RequestInitChain) abci.ResponseInitChain {
	bridgemodule.InitGenesis(ctx)
	return abci.ResponseInitChain{}
}

func noopTxDecoder(_ []byte) (sdk.Tx, error) {
	return nil, errors.New("tx decoding not configured for QubeticsApp skeleton")
}

type noopBankKeeper struct{}

func (noopBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

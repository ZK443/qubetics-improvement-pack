package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bridgekeeper "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
)

type QubeticsApp struct {
	*baseapp.BaseApp
	BridgeKeeper bridgekeeper.Keeper
}

func NewQubeticsApp() *QubeticsApp {
	app := &QubeticsApp{
		BaseApp: baseapp.NewBaseApp("qubetics", nil, store.NewCommitMultiStore(nil)),
	}

	// Инициализация хранилища для bridge-модуля
	keyBridge := sdk.NewKVStoreKey("bridge")
	ctx := sdk.NewContext(store.NewCommitMultiStore(nil), abci.Header{}, false, nil)

	app.BridgeKeeper = bridgekeeper.Keeper{
		StoreService: prefix.NewStore(ctx.KVStore(keyBridge), []byte("bridge/")),
	}

	return app
}

func (app *QubeticsApp) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (app *QubeticsApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (app *QubeticsApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	app.BridgeKeeper.InitGenesis(ctx)
	return abci.ResponseInitChain{}
}

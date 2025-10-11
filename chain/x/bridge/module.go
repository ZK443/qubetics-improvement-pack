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
func (AppModul

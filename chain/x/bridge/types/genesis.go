package types

// GenesisState — минимальный набор для инициализации модуля bridge.
type GenesisState struct {
	GlobalPause bool            `json:"global_pause" yaml:"global_pause"`
	Params      Params          `json:"params"       yaml:"params"`
	ACL         map[string]bool `json:"acl"          yaml:"acl"`
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		GlobalPause: false,
		Params:      DefaultParams(),
		ACL:         map[string]bool{},
	}
}

func (g *GenesisState) Validate() error {
	// синхронизировать флаг паузы с Params
	g.Params.GlobalPause = g.GlobalPause
	return g.Params.Validate()
}

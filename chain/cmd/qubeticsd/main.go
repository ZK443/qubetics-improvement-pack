//go:build cosmos

package main

import (
	"os"

	"github.com/spf13/cobra"

	bridgecli "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/client/cli"
)

func main() {
	root := &cobra.Command{
		Use:   "qubeticsd",
		Short: "Qubetics node/cli (skeleton)",
	}

	// tx namespace
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}
	txCmd.AddCommand(bridgecli.NewTxCmd())
	root.AddCommand(txCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

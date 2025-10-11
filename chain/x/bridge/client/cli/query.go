//go:build cosmos

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	brtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge",
		Short: "Bridge queries",
	}
	cmd.AddCommand(
		NewQueryParamsCmd(),
		NewQueryStatusCmd(),
		NewQueryNonceCmd(),
		NewQueryIsAllowedCmd(),
	)
	return cmd
}

// bridge query params
func NewQueryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Show bridge parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			q := brtypes.NewQueryClient(ctx)
			resp, err := q.Params(cmd.Context(), &brtypes.ParamsRequest{})
			if err != nil {
				return err
			}
			return ctx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// bridge query status <message-id>
func NewQueryStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [message-id]",
		Short: "Get message status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			q := brtypes.NewQueryClient(ctx)
			resp, err := q.Status(cmd.Context(), &brtypes.StatusRequest{MessageId: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// bridge query nonce <sender>
func NewQueryNonceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nonce [sender]",
		Short: "Get current sender nonce",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			q := brtypes.NewQueryClient(ctx)
			resp, err := q.Nonce(cmd.Context(), &brtypes.NonceRequest{Sender: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// bridge query is-allowed <address>
func NewQueryIsAllowedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-allowed [address]",
		Short: "Check ACL permission for address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			q := brtypes.NewQueryClient(ctx)
			resp, err := q.IsAllowed(cmd.Context(), &brtypes.IsAllowedRequest{Addr: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// helper (не используется напрямую, но оставлен на будущее)
func must(err error) {
	if err != nil {
		panic(fmt.Errorf("bridge query: %w", err))
	}
}

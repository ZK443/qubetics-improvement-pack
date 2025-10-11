//go:build cosmos

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktx "github.com/cosmos/cosmos-sdk/client/tx"

	brtypes "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge",
		Short: "Bridge transactions",
	}
	cmd.AddCommand(
		NewVerifyProofCmd(),
		NewExecuteCmd(),
	)
	return cmd
}

func NewVerifyProofCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-proof [message-id] [proof-id] [verifier]",
		Short: "Submit a proof for verification",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &brtypes.MsgVerifyProof{
				Signer:     clientCtx.GetFromAddress().String(),
				Proof:      &brtypes.Proof{ProofId: args[1], Data: []byte{}},
				MessageId:  args[0],
				Verifier:   args[2],
				// ExpectedDigest опционален
			}

			if err := validateMsgVerifyProof(msg); err != nil {
				return err
			}

			return sdktx.GenerateOrBroadcastTxCLI(clientCtx, sdktx.NewFactoryCLI(clientCtx, cmd.Flags()), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewExecuteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [message-id] [route: token-transfer|contract-call]",
		Short: "Execute verified message",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			route, err := parseExecRoute(args[1])
			if err != nil {
				return err
			}

			msg := &brtypes.MsgExecute{
				Executor:  clientCtx.GetFromAddress().String(),
				MessageId: args[0],
				Route:     route,
				// Amount строкой, если позже потребуется лимит по сумме
			}

			if err := validateMsgExecute(msg); err != nil {
				return err
			}

			return sdktx.GenerateOrBroadcastTxCLI(clientCtx, sdktx.NewFactoryCLI(clientCtx, cmd.Flags()), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// --- helpers ---

func parseExecRoute(s string) (brtypes.ExecRoute, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "token-transfer", "token", "transfer":
		return brtypes.ExecRoute_EXEC_ROUTE_TOKEN_TRANSFER, nil
	case "contract-call", "contract":
		return brtypes.ExecRoute_EXEC_ROUTE_CONTRACT_CALL, nil
	default:
		return brtypes.ExecRoute_EXEC_ROUTE_UNSPECIFIED, fmt.Errorf("unknown route: %s", s)
	}
}

func validateMsgVerifyProof(m *brtypes.MsgVerifyProof) error {
	if m == nil || m.Signer == "" || m.MessageId == "" || m.Proof == nil || m.Proof.ProofId == "" {
		return fmt.Errorf("missing required fields")
	}
	return nil
}

func validateMsgExecute(m *brtypes.MsgExecute) error {
	if m == nil || m.Executor == "" || m.MessageId == "" {
		return fmt.Errorf("missing required fields")
	}
	if m.Route == brtypes.ExecRoute_EXEC_ROUTE_UNSPECIFIED {
		return fmt.Errorf("route must be specified")
	}
	return nil
}

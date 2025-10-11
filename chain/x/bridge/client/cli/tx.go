package cli

import (
	"fmt"

	"github.com/spf13/cobra"
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
			// Заглушка: здесь формируется Tx с MsgVerifyProof и отправляется в ноду.
			// В skeleton-приложении выводится инструкция.
			fmt.Fprintf(cmd.OutOrStdout(), "prepare tx: verify-proof msg_id=%s proof_id=%s verifier=%s\n", args[0], args[1], args[2])
			return nil
		},
	}
	return cmd
}

func NewExecuteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute [message-id] [route: token-transfer|contract-call]",
		Short: "Execute verified message",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "prepare tx: execute msg_id=%s route=%s\n", args[0], args[1])
			return nil
		},
	}
	return cmd
}

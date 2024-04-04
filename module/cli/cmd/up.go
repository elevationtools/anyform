package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run the up DAG.",
	SilenceUsage: true,
	SilenceErrors: true,
	// Long: "",
	RunE: func(cmd *cobra.Command, args []string) error {
    orc := Must(AnyformSingleton().NewOrchestrator())
		return orc.Up(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

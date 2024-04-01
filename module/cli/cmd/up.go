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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

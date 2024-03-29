package cmd

import (
  "context"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove the ANYFORM_GENFILES directory",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {
    Must1(Must(AnyformSingleton().NewOrchestrator()).Clean(context.Background()))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}


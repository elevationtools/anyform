package cmd

import (
	"github.com/spf13/cobra"
	. "github.com/elevationtools/anyform/common/util"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run the up DAG.",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {
    orc := Must(Singleton().NewOrchestrator())
		Must1(orc.Up())
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

package cmd

import (
	"github.com/spf13/cobra"
	. "github.com/elevationtools/anyform/common/util"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run the down DAG",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {
    Must(Singleton().NewOrchestrator()).Down()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


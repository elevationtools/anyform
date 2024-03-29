package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	commonutil "github.com/elevationtools/anyform/common/util"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Print the orchestrator spec",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {
    orc := Must(AnyformSingleton().NewOrchestrator())
		fmt.Printf("Orchestrator Spec: ")
		fmt.Println(Must(commonutil.ToJSONString(orc.Spec)))
	},
}

func init() {
	rootCmd.AddCommand(specCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// specCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// specCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


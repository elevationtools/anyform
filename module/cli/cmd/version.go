package cmd

import (
  "fmt"

	"github.com/spf13/cobra"
)

const (
  AnyformVersion = "0.0.0-unsetgocode"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Run: func(cmd *cobra.Command, args []string) {
    fmt.Printf("Anyform %v\n", AnyformVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

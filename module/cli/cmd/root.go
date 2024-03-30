package cmd

import (
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anyform",
	Short: "Infrastructure-as-code tool",
	Long: `Infrastructure-as-code tool

Anyform is meant to fill gaps in related tools like Helm, Jsonnet, Terraform,
etc.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevelStr := Must(cmd.Flags().GetString("loglevel"))
		slog.SetLogLoggerLevel((func () slog.Level {
			switch logLevelStr {
				case "all": return -1e6
				case "debug": return slog.LevelDebug
				case "info": return slog.LevelInfo
				case "warn": return slog.LevelWarn
				case "error": return slog.LevelError
				case "off": return 1e6
			}
			log.Fatalf("Invalid log level: %v", logLevelStr)
			return 0
		})())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().String("loglevel", "warn",
			"Display log levels at or above this value. Values: debug, info, warn, error")
}

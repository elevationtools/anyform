
package main

import (
  "context"
  "fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

  anyform "github.com/elevationtools/anyform/module/lib"
	commonutil "github.com/elevationtools/anyform/module/common/util"
)

func InitRoot(orc *anyform.Orchestrator, orcErr error) *cobra.Command {
  var rootCmd = &cobra.Command{
    Use:   "anyform",
    Short: "Infrastructure-as-code tool",
    Long: `Infrastructure-as-code tool

  Anyform is meant to fill gaps in related tools like Helm, Jsonnet, Terraform,
  etc.`,

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

	rootCmd.PersistentFlags().String("loglevel", "warn",
			"Display log levels at or above this value. Values: debug, info, warn, error")

  InitMark(orc, orcErr, rootCmd)

  rootCmd.AddCommand(&cobra.Command{
    Use: "up",
    Long: "Run the up DAG.",
    SilenceUsage: true,
    SilenceErrors: true,
    RunE: RunEWrapper(orc, orcErr, func(cmd *cobra.Command, args []string) error {
      return orc.Up(context.Background())
    }),
  })

  rootCmd.AddCommand(&cobra.Command{
    Use: "down",
    Long: "Run the down DAG",
    SilenceUsage: true,
    SilenceErrors: true,
    RunE: RunEWrapper(orc, orcErr, func(cmd *cobra.Command, args []string) error {
      return orc.Down(context.Background())
    }),
  })

	rootCmd.AddCommand(&cobra.Command{
    Use: "clean",
    Long: "Remove the ANYFORM_GENFILES directory",
    RunE: RunEWrapper(orc, orcErr, func(cmd *cobra.Command, args []string) error {
      return orc.Clean(context.Background())
    }),
  })

  rootCmd.AddCommand(&cobra.Command{
    Use: "spec",
    Long: "Print the orchestrator spec",
    RunE: RunEWrapper(orc, orcErr, func(cmd *cobra.Command, args []string) error {
      fmt.Printf("Orchestrator Spec: ")
      fmt.Println(Must(commonutil.ToJSONString(orc.Spec)))
      return nil
    }),
  })

  rootCmd.AddCommand(&cobra.Command{
    Use:   "version",
    Short: "Print the version.",
    RunE: func(cmd *cobra.Command, args []string) error {
      fmt.Printf("Anyform %v\n", AnyformVersion)
      return nil
    },
  })

  return rootCmd
}

var (
  AnyformVersion = "0.0.0-unsetgocode"
)



package main

import (
  "fmt"

	"github.com/spf13/cobra"

  anyform "github.com/elevationtools/anyform/module/lib"
)

func InitMark(orc *anyform.Orchestrator, orcErr error, rootCmd *cobra.Command) {
  var markCmd = &cobra.Command{
    Use: "mark",
    Long: "Force assuming a stage is in a given state.",
  }

  markCommon := func(state string) *cobra.Command {
    return &cobra.Command{
      Use: state,
      Long: fmt.Sprintf("Force assuming the stage is %v.", state),
      SilenceUsage: true,
      SilenceErrors: false,
      Args: cobra.MinimumNArgs(1),
      RunE: RunEWrapper(orc, orcErr, func(cmd *cobra.Command, args []string) error {
        stages := []*anyform.Stage{}
        for _, stageName := range args {
          stage, err := orc.GetStage(stageName)
          if err != nil { return err }
          stages = append(stages, stage)
        }
        errors := map[string]error{}
        for _, stage := range stages {
          err := stage.Mark(state)
          if err != nil { errors[stage.Name] = err }
        }
        if len(errors) == 0 {
          return nil
        }

        // TODO(ux): this will be very ugly.
        msg := "stages with errors:\n"
        for stageName, err := range errors {
          msg += fmt.Sprintf("\n  %v: %v\n", stageName, err)
        }
        return fmt.Errorf(msg)
      }),
    }
  }

  rootCmd.AddCommand(markCmd)
  markCmd.AddCommand(markCommon("up"))
  markCmd.AddCommand(markCommon("down"))
}


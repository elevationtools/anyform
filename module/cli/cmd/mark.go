package cmd

import (
  "fmt"

	"github.com/spf13/cobra"

  anyform "github.com/elevationtools/anyform/module/lib"
)

var markCmd = &cobra.Command{
	Use: "mark",
	Long: "Force assuming the stage is up or down.",
}

var markUpCmd = markCommon("up")
var markDownCmd = markCommon("down")

func markCommon(state string) *cobra.Command {
  return &cobra.Command{
    Use: state,
    Long: fmt.Sprintf("Force assuming the stage is %v.", state),
    SilenceUsage: true,
    SilenceErrors: false,
    Args: cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      orc := Must(AnyformSingleton().NewOrchestrator())
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
    },
  }
}

func init() {
	rootCmd.AddCommand(markCmd)
	markCmd.AddCommand(markUpCmd)
	markCmd.AddCommand(markDownCmd)
}


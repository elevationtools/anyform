
package main

import (
	"github.com/spf13/cobra"

  anyform "github.com/elevationtools/anyform/module/lib"
)

func Must[T any](result T, err error) T {
  if err != nil { panic(err) }
  return result
}

func Must1(err error) {
  if err != nil { panic(err) }
}

type RunEFunc func (cmd *cobra.Command, args []string) error

func RunEWrapper(orc *anyform.Orchestrator, orcErr error, runEFunc RunEFunc) RunEFunc {
  return func(cmd *cobra.Command, args []string) error {
    if orcErr != nil { return orcErr }
    return runEFunc(cmd, args)
  }
}
 

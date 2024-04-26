
package main

import (
	"os"

  anyform "github.com/elevationtools/anyform/module/lib"
)

func main() {
  orc, err := anyform.NewDefaultAnyform().NewOrchestrator()
  rootCmd := InitRoot(orc, err)
	err = rootCmd.Execute()
	if err != nil { os.Exit(1) }
}


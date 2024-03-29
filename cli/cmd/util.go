
package cmd

import (
)

func Must[T any](result T, err error) T {
  if err != nil { panic(err) }
  return result
}

func Must1(err error) {
  if err != nil { panic(err) }
}

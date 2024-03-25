
package anyform

import (
  "fmt"
)

type Anyform struct {
}

func New() *Anyform {
  return &Anyform{}
}

func (af* Anyform) Up() {
  fmt.Println("lib-anyform Up!")
}

func (af* Anyform) Down() {
  fmt.Println("lib-anyform Down!")
}


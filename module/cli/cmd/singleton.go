
package cmd

import (
  "sync"

  anyform "github.com/elevationtools/anyform/module/lib"
)

var singleton *anyform.Anyform;
var singletonOnce sync.Once;

func AnyformSingleton() *anyform.Anyform {
  singletonOnce.Do(func () {
    singleton = anyform.NewDefaultAnyform()
  })
  return singleton
}

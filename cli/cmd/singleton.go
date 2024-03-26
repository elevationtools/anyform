
package cmd

import (
  "sync"

  anyform "github.com/elevationtools/anyform/lib"
)

var singleton *anyform.Anyform;
var singletonOnce sync.Once;

func Singleton() *anyform.Anyform {
  singletonOnce.Do(func () {
    singleton = anyform.NewAnyform()
  })
  return singleton
}



package dag

import (
  "context"
)

type VertexState = int
const (
  VertexStateUnknown VertexState = 0
  VertexStateNotStarted VertexState = 1
  VertexStateRunning VertexState = 2
  VertexStateSuccess VertexState = 3
  VertexStateFailure VertexState = 4
)

type VertexRunFunc = func(ctx context.Context) error

type Dag interface {
  AddVertex(name string, childNames []string, runFunc VertexRunFunc) error
  Run(ctx context.Context) error
}

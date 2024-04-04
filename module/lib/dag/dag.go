
package dag

import (
  "context"
  "fmt"
)

type VertexState int
const (
  VertexStateUnknown VertexState = 0
  VertexStateNotStarted VertexState = 1
  VertexStateRunning VertexState = 2
  VertexStateSuccess VertexState = 3
  VertexStateFailure VertexState = 4
)

func (vs VertexState) ToString() string {
  switch (vs) {
    case VertexStateUnknown: return "VertexStateUnknown"
    case VertexStateNotStarted: return "VertexStateNotStarted"
    case VertexStateRunning: return "VertexStateRunning"
    case VertexStateSuccess: return "VertexStateSuccess"
    case VertexStateFailure: return "VertexStateFailure"
    default: panic(fmt.Sprintf("Unknown VertexState %d", vs))
  }
}

type VertexRunFunc = func(ctx context.Context) error

type Dag interface {
  AddVertex(name string, childNames []string, runFunc VertexRunFunc) error
	// If `parallel` is false, only run one stage at a time, otherwise use maximum
	// parallelism.
  Run(ctx context.Context, parallel bool) error
	// TODO: add ability to see the errors, perhaps return []*Vertex instead (would
  // require adding Vertex to the interface).
	FailedVerticies() []string
}

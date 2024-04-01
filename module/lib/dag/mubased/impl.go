
package dag_mutex_based

import (
  "container/list"
	"context"
  "fmt"
	"maps"
	"strings"
  "sync"

  daglib "github.com/elevationtools/anyform/module/lib/dag"
)

// Vertex //////////////////////////////////////////////////////////////////////

type Vertex struct {
  name string
  runFunc daglib.VertexRunFunc
  state daglib.VertexState  // Synchronized by Dag.mu
  parentNames []string
  children []*Vertex
  parents []*Vertex
  notStartedElement *list.Element
  runningElement *list.Element
}

func newVertex(name string, parentNames []string, runFunc daglib.VertexRunFunc) *Vertex {
  return &Vertex{
    name: name,
    runFunc: runFunc,
		// synchronized by Dag.mu
    state: daglib.VertexStateUnknown,
    parentNames: parentNames,
    children: []*Vertex{},
    parents: []*Vertex{},
  }
}

func (v *Vertex) isReady() bool {
  for _, parent := range v.parents {
    if parent.state != daglib.VertexStateSuccess {
      return false
    }
  }
  return true
}

// Dag /////////////////////////////////////////////////////////////////////////////

// Methods are thread compatible unless otherwise specified.
type Dag struct {
  all map[string]*Vertex
  roots []*Vertex

	// The following are synchronized by mu
  notStarted list.List
  running list.List
  succeeded []*Vertex
  failed []*Vertex

  done chan bool
  mu sync.Mutex
}

func NewDag() daglib.Dag {
  type vl = []*Vertex
  return &Dag {
    all: map[string]*Vertex{},
		roots: vl{},
    succeeded: vl{},
    failed: vl{},
		done: make(chan bool),
  }
}

func (dag *Dag) AddVertex(name string, parentNames []string,
                          runFunc daglib.VertexRunFunc) error {
  // Prevent duplicates
  _, found := dag.all[name]
  if found {
    return fmt.Errorf("Dag.AddVertex: duplicate vertex added: %v", name)
  }

  v := newVertex(name, parentNames, runFunc)
  dag.all[name] = v
  return nil
}

func (dag *Dag) Run(ctx context.Context) error {
	// Create the DAG
  for _, v := range dag.all {
		// dag.mu not needed because no concurrency has started yet.
		v.state = daglib.VertexStateNotStarted
    v.notStartedElement = dag.notStarted.PushBack(v)

    if len(v.parentNames) == 0 {
      dag.roots = append(dag.roots, v)
    }

    // Populate Vertex.children and the children's Vertex.parents.
    for _, parentName := range v.parentNames {
      parent, found := dag.all[parentName]
      if !found {
        return fmt.Errorf("Vertex with missing parent: vertex=%v missing_parent=%v",
                          v.name, parentName)
      }

      v.parents = append(v.parents, parent)
      parent.children = append(parent.children, v)
    }
  }

	// Make sure it's acyclic
	for _, v := range dag.roots {
		path := dag.getCycle(v, map[*Vertex]bool{}, []*Vertex{})
		if path != nil {
			names := []string{}
			for _, j := range path {
				names = append(names, j.name)
			}
			return fmt.Errorf("Cycle detected: %v", strings.Join(names, ", "))
		}
	}

	// Start all the roots.
  dag.mu.Lock()
  for _, v := range dag.roots {
    dag.transitionNotStartedToRunning_prelock(v, ctx)
  }
  dag.mu.Unlock()

	select {
		case <- ctx.Done(): return ctx.Err()
		case <- dag.done:
	}

	if len(dag.failed) > 0 {
		return fmt.Errorf("verticies failed: %v", dag.FailedVerticies())
	}
	return nil
}

// Returns empty if no failures.
func (dag *Dag) FailedVerticies() []string {
	names := []string{}
	for _, v := range dag.failed {
		names = append(names, v.name)
	}
	return names
}

// Returns nil if there is no cycle.
// Others, returns the path causing the cycle starting from the root.
func (dag *Dag) getCycle(v *Vertex, seen map[*Vertex]bool, path []*Vertex) []*Vertex {
	path = append(path, v)
	if seen[v] {
		return path
	}
	seen = maps.Clone(seen)
	seen[v] = true
	for _, child := range v.children {
		newPath := dag.getCycle(child, seen, path)
		if newPath != nil {
			return newPath
		}
	}
	return nil
}

func (dag *Dag) transitionNotStartedToRunning_prelock(v *Vertex,
																									    ctx context.Context) {
  if v.notStartedElement == nil { panic("unexpected") }
  if v.runningElement != nil { panic("unexpected") }

	// Don't start if we already got cancelled.
	if ctx.Err() != nil { return }

  v.state = daglib.VertexStateRunning
  v.runningElement = dag.running.PushBack(v)
  dag.notStarted.Remove(v.notStartedElement)
  v.notStartedElement = nil

  go func() {
		if v.notStartedElement != nil { panic("unexpected") }
		if v.runningElement == nil { panic("unexpected") }
		if v.state != daglib.VertexStateRunning { panic("unexpected") }

		err := v.runFunc(ctx)

		dag.mu.Lock()
		defer dag.mu.Unlock()
		childrenStarted := 0  // just for sanity checking.
		if err != nil {
			dag.transitionOutOfRunning_prelock(v, false)
		} else {
			dag.transitionOutOfRunning_prelock(v, true)
			for _, child := range v.children {
				if child.state == daglib.VertexStateNotStarted && child.isReady() {
					childrenStarted++
					dag.transitionNotStartedToRunning_prelock(child, ctx)
				}
			}
		}

		// Nothing to start
		if dag.running.Len() == 0 {
			if childrenStarted > 0 { panic("unexpected") }
			close(dag.done)
		}
	}()
}

func (dag *Dag) transitionOutOfRunning_prelock(v *Vertex, success bool) {
  if success {
    v.state = daglib.VertexStateSuccess
    dag.succeeded = append(dag.succeeded, v)
  } else {
    v.state = daglib.VertexStateFailure
    dag.failed = append(dag.failed, v)
  }
  if v.runningElement == nil { panic("unexpected") }
  dag.running.Remove(v.runningElement)
  v.runningElement = nil
}

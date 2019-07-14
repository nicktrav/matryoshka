package graph

import (
	"testing"

	"github.com/nicktrav/matryoshka/pkg/lang"
)

// newGraph returns a new Dependency graph.
//
// The graph is of the form:
//
//            foo
//          /     \
//        bar     bam
//        / \    /
//     baz   boom
//
func newGraph() *DependencyGraph {
	graph := NewDependencyGraph()
	graph.Construct([]*lang.Dep{fooRawDep, barRawDep, bazRawDep, bamRawDep, boomRawDep})
	return graph
}

func TestDepthFirstWalker_Walk(t *testing.T) {
	graph := newGraph()

	tracker := newTracker()
	walker := NewWalker(tracker)

	err := walker.Walk(graph, "foo")
	if err != nil {
		t.Fatalf("got error: %+v", err)
	}

	walk := tracker.depsVisited
	if len(walk) != 5 {
		t.Fatalf("wanted size of walk to be 5; got %d", len(walk))
	}

	// expect a post-order traversal of the dep graph
	expectedOrder := []string{"baz", "boom", "bar", "bam", "foo"}

	for i, want := range expectedOrder {
		got := walk[i]
		if got.Name != want {
			t.Errorf("wanted visited dep #%d to be %s; got %+v", i, want, got)
		}
	}
}

// pathTracker is a DepVisitor that maintains an ordered list of Deps visited.
type pathTracker struct {

	// depsVisited is a slice containing pointers to Deps visited by the
	// walker.
	depsVisited []*Dependency
}

// newTracker returns a new pathTracker
func newTracker() *pathTracker {
	return &pathTracker{depsVisited: []*Dependency{}}
}

// Visit adds the dep to the deps visited
func (t *pathTracker) Visit(dep *Dependency) error {
	t.depsVisited = append(t.depsVisited, dep)
	return nil
}

func (t *pathTracker) PreVisit(dep *Dependency) {
}

func (t *pathTracker) PostVisit(dep *Dependency) {
}

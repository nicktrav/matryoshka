package graph

import (
	"testing"

	"github.com/nicktrav/matryoshka/pkg/lang"
)

func TestDependencyGraph_Construct(t *testing.T) {
	g := NewDependencyGraph()
	g.Construct([]*lang.Dep{bamRawDep, boomRawDep})

	// graph contains two deps
	if len(g.Deps()) != 2 {
		t.Fatalf("wanted 2 entries in the graph; got %+v", g.Deps())
	}

	// boom is well formed

	boomDep := g.depMap[boomRawDep.Name]
	if boomDep.Name != boomRawDep.Name {
		t.Errorf("wanted name %s, got %s", boomDep.Name, boomDep.Name)
	}

	if len(boomDep.Dependencies) != 0 {
		t.Errorf("wanted no deps for boom; got %d", len(boomDep.Dependencies))
	}

	// bam is well formed

	bamDep := g.depMap[bamRawDep.Name]
	if bamDep.Name != bamRawDep.Name {
		t.Errorf("wanted name %s, got %s", bamRawDep.Name, bamDep.Name)
	}

	if len(bamDep.Dependencies) != 1 {
		t.Errorf("wanted requirements slice to be of size 1; got %d", len(bamDep.Dependencies))
	}

	if bamDep.Dependencies[0] != boomDep {
		t.Errorf("wanted bam's only dep to be boom; got %v", bamDep.Dependencies[0])
	}

	// converter map should contain only the deps foo and bar

	if len(g.depMap) != 2 {
		t.Errorf("wanted map to be of size 2; got %d", len(g.depMap))
	}

	assertMapContainsDep(t, g, "bam", bamDep)
	assertMapContainsDep(t, g, "boom", boomDep)
}

func assertMapContainsDep(t *testing.T, g *DependencyGraph, depName string, wantedDep *Dependency) {
	dep, found := g.depMap[depName]

	if !found {
		t.Errorf("wanted map to contain key '%s'; got %+v", depName, g.depMap)
	}

	if dep != wantedDep {
		t.Errorf("wanted '%s' key to be dep %s; was %+v", depName, wantedDep.Name, dep)
	}
}

func TestDependencyGraph_Deps(t *testing.T) {
	g := NewDependencyGraph()
	g.depMap["foo"] = g.makeDep(fooRawDep)

	if len(g.Deps()) != 5 {
		t.Fatalf("wanted deps have length 5; got %d", len(g.Deps()))
	}

	wanted := []string{"foo", "bar", "baz", "bam", "boom"}
	for _, want := range wanted {
		contains(t, want, g.Deps())
	}
}

func TestDependencyGraph_Get_DepMissing(t *testing.T) {
	g := NewDependencyGraph()

	dep := g.Get("foo")
	if dep != nil {

	}
}

func TestDependencyGraph_Get_DepPresent(t *testing.T) {
	g := NewDependencyGraph()
	g.depMap[fooRawDep.Name] = g.makeDep(fooRawDep)

	dep := g.Get(fooRawDep.Name)
	if dep == nil {
		t.Fatalf("wanted dep to be non-nil; got nil")
	}

	if dep.Name != fooRawDep.Name {
		t.Errorf("wanted dep to be %s; got %s", fooRawDep.Name, dep.Name)
	}
}

// contains checks that the dep with the given name is in the list of deps.
func contains(t *testing.T, wanted string, deps []*Dependency) {
	for _, dep := range deps {
		if dep.Name == wanted {
			return
		}
	}
	t.Errorf("wanted list of deps to contain %s", wanted)
}

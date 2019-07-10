package pkg

import (
	"testing"

	"github.com/nicktrav/matryoshka/pkg/lang"
)

func TestDependencyGraph_Construct(t *testing.T) {
	g := NewDependencyGraph()

	fooRawDep := &lang.Dep{
		Name:         "foo",
		Requirements: []*lang.Dep{},
		MetCommands:  []*lang.ShellCmd{},
		MeetCommands: []*lang.ShellCmd{},
	}

	barRawDep := &lang.Dep{
		Name:         "bar",
		Requirements: []*lang.Dep{fooRawDep},
		MetCommands:  []*lang.ShellCmd{},
		MeetCommands: []*lang.ShellCmd{},
	}

	// foo is well formed

	fooDep := g.makeDep(fooRawDep)

	if fooDep.Name != fooRawDep.Name {
		t.Errorf("wanted name %s, got %s", fooRawDep.Name, fooDep.Name)
	}

	if len(fooDep.Dependencies) != 0 {
		t.Errorf("wanted no deps for foo; got %d", len(fooDep.Dependencies))
	}

	// bar is well formed

	barDep := g.makeDep(barRawDep)

	if barDep.Name != barRawDep.Name {
		t.Errorf("wanted name %s, got %s", barDep.Name, barDep.Name)
	}

	if len(barDep.Dependencies) != 1 {
		t.Errorf("wanted requirements slice to be of size 1; got %d", len(barDep.Dependencies))
	}

	if barDep.Dependencies[0] != fooDep {
		t.Errorf("wanted bar's only dep to be foo; got %v", barDep.Dependencies[0])
	}

	// converter map should contain only the deps foo and bar

	if len(g.depMap) != 2 {
		t.Errorf("wanted map to be of size 2; got %d", len(g.depMap))
	}

	assertMapContainsDep(t, g, "foo", fooDep)
	assertMapContainsDep(t, g, "bar", barDep)
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

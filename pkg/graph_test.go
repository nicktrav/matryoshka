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
		t.Errorf("expected no deps for foo; got %d", len(fooDep.Dependencies))
	}

	// bar is well formed

	barDep := g.makeDep(barRawDep)

	if barDep.Name != barRawDep.Name {
		t.Errorf("wanted name %s, got %s", barDep.Name, barDep.Name)
	}

	if len(barDep.Dependencies) != 1 {
		t.Errorf("expected requirements slice to be of size 1; got %d", len(barDep.Dependencies))
	}

	if barDep.Dependencies[0] != fooDep {
		t.Errorf("expected bar's only dep to be foo; was %v", barDep.Dependencies[0])
	}

	// converter map should contain only the deps foo and bar

	if len(g.depMap) != 2 {
		t.Errorf("expected map to be of size 2; was %d", len(g.depMap))
	}

	assertMapContainsDep(t, g, "foo", fooDep)
	assertMapContainsDep(t, g, "bar", barDep)
}

func assertMapContainsDep(t *testing.T, g *DependencyGraph, depName string, expectedDep *Dependency) {
	dep, found := g.depMap[depName]

	if !found {
		t.Errorf("expected map to contain key '%s'", depName)
	}

	if dep != expectedDep {
		t.Errorf("expected '%s' key to be dep %s; was %v", depName, expectedDep.Name, dep)
	}
}

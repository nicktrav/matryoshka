package graph

import (
	"github.com/nicktrav/matryoshka/pkg/actions"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

// DependencyGraph is a directed, acyclic graph of Dependencies. Nodes are
// Dependencies and edges represent a dependency on another Dependency.
type DependencyGraph struct {

	// depMap is a mapping from dependency name to Dependency.
	depMap map[string]*Dependency
}

// NewDependencyGraph returns a pointer to a new DependencyGraph.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		depMap: make(map[string]*Dependency),
	}
}

// Construct takes a slice of mappings of dep names to Values and populates the
// DependencyGraph.
//
// The graph is constructed by flattening the mappings of names to Values,
// filtering only Dep types, constructing a new Dependency and placing it in
// the dep map.
func (g *DependencyGraph) Construct(deps []*lang.Dep) {
	for _, dep := range deps {
		// exclude any deps that aren't enabled
		if !dep.Enable {
			continue
		}

		if _, ok := g.depMap[dep.Name]; ok {
			// TODO(nickt) this implies there's a duplicate dep, and we should warn
			continue
		}
		g.makeDep(dep)
	}
}

// makeDep translates a Dep into a new Dependency, using a cached value if a
// Dependency with the same name is already present in the dep map.
func (g *DependencyGraph) makeDep(rawDep *lang.Dep) *Dependency {
	// deps that aren't enabled can still end up in the graph if referenced
	// directly from the requirements block of another dep, and should not be
	// added to the graph. Return nil as a sentinel value.
	if !rawDep.Enable {
		return nil
	}

	// if dep is already in the map, return it
	dep, found := g.depMap[rawDep.Name]
	if found {
		return dep
	}

	// else, construct the dependency
	dep = &Dependency{
		Name:        rawDep.Name,
		MetActions:  convertCommands(rawDep.MetCommands),
		MeetActions: convertCommands(rawDep.MeetCommands),
	}

	// for each requirement, recurse
	var requirements []*Dependency
	for _, req := range rawDep.Requirements {
		reqDep := g.makeDep(req)
		if reqDep != nil {
			requirements = append(requirements, reqDep)
		}
	}

	// add the generated Deps into the requirements list
	dep.Dependencies = requirements

	// and place this dep into the map
	g.depMap[rawDep.Name] = dep

	return dep
}

// Get returns the Dependency with the given name from the graph.
func (g *DependencyGraph) Get(name string) *Dependency {
	dep, found := g.depMap[name]
	if !found {
		return nil
	}
	return dep
}

// Deps returns the a slice of all Dependencies in the DependencyGraph.
func (g *DependencyGraph) Deps() []*Dependency {
	var deps []*Dependency
	for _, val := range g.depMap {
		deps = append(deps, val)
	}
	return deps
}

// convertCommands takes a slice of pointers to ShellCmds and converts them
// into a slice of Actions.
func convertCommands(shellCommands []*lang.ShellCmd) []actions.Action {
	var as []actions.Action
	for _, command := range shellCommands {
		shell := actions.NewShellCommandAction(command)
		as = append(as, shell)
	}
	return as
}

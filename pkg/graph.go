package pkg

import (
	"fmt"

	"go.starlark.net/starlark"

	"github.com/nicktrav/matryoshka/pkg/actions"
	"github.com/nicktrav/matryoshka/pkg/lang"
)

// DependencyGraph is a directed, acyclic graph of Dependencies. Nodes are
// Dependencies and edges represent a dependency on another Dependency.
type DependencyGraph struct {

	// depMap is a mapping from dependency name to Dependency.
	depMap map[string]*Dependency
}

// A nodeVisitor iterates over nodes in the DependencyGraph
type NodeVisitor interface {

	// Visit is the main action to take when visiting a node in the graph.
	Visit(dep *Dependency)

	// PreVisit is a hook called before the Visit function is called.
	PreVisit(dep *Dependency)

	// PostVisit is a hook called after the Visit function is called.
	PostVisit(dep *Dependency)

	// Errors returns the accumulation of errors seen while visiting nodes in
	// the DependencyGraph.
	Errors() []error
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
func (g *DependencyGraph) Construct(globals []starlark.StringDict) {
	// for each dep in the flattened list of deps
	for _, dict := range globals {
		for _, global := range dict.Keys() {
			value := dict[global]

			// skip any global that's not a Dep
			rawDep, ok := value.(*lang.Dep)
			if !ok {
				continue
			}

			// check if the dep is already in the map
			if _, found := g.depMap[global]; found {
				// TODO(nickt) this implies there's a duplicate dep, and we should warn
				continue
			}

			// else, construct the dep
			g.makeDep(rawDep)
		}
	}
}

// makeDep translates a Dep into a new Dependency, using a cached value if a
// Dependency with the same name is already present in the dep map.
func (g *DependencyGraph) makeDep(rawDep *lang.Dep) *Dependency {
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
		requirements = append(requirements, g.makeDep(req))
	}

	// add the generated Deps into the requirements list
	dep.Dependencies = requirements

	// and place this dep into the map
	g.depMap[rawDep.Name] = dep

	return dep
}

// Walk starts a depth-first post-order traversal of the graph, starting at
// the Dependency with the given name.
func (g *DependencyGraph) Walk(name string, visitor NodeVisitor) error {
	start := g.Get(name)
	if start == nil {
		return fmt.Errorf("node %s not found", name)
	}
	return g.visit(start, visitor, make(map[string]int))
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

// visit uses a nodeVisitor to attempt to visit the given Dependency,
// recursively visiting the Dependency's own dependencies by calling visit on
// them.
// TODO(nickt): detect cycles
// TODO(nickt): move visit cache into the nodeVisitor
// TODO(nickt): consider moving visit into the actual visitor impl
func (g *DependencyGraph) visit(dep *Dependency, visitor NodeVisitor, visitCache map[string]int) error {
	visitor.PreVisit(dep)
	defer visitor.PostVisit(dep)

	if dep == nil {
		return fmt.Errorf("dependency not found in graph")
	}

	// if we've visited this node, no need to visit again
	count, visited := visitCache[dep.Name]
	if visited {
		return nil
	}

	// however, if the count on this node is greater than zero, we're already visiting
	// this node higher up in the graph
	if count > 0 {
		return fmt.Errorf("detected cycle")
	}

	// and add this to the list of visited nodes
	visitCache[dep.Name]++

	// and visit all of the dependent nodes
	for _, dep := range dep.Dependencies {
		dependency := g.Get(dep.Name)
		err := g.visit(dependency, visitor, visitCache)
		if err != nil {
			return fmt.Errorf("error visiting dependency '%s': %s", dep.Name, err)
		}
	}

	// take action on this current node
	visitor.Visit(dep)

	// "un-visit" this dep from the list of nodes
	visitCache[dep.Name]--

	return nil
}

// convertCommands takes a slice of pointers to ShellCmds and converts them
// into a slice of Actions.
func convertCommands(shellCommands []*lang.ShellCmd) []actions.Action {
	var as []actions.Action
	for _, command := range shellCommands {
		shell := actions.NewShellCommandAction(command.Command)
		as = append(as, shell)
	}
	return as
}

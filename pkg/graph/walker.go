package graph

import (
	"fmt"
)

// Walker walks a DependencyGraph, visiting the deps in a given order.
type Walker interface {

	// Walk traverses the given DependencyGraph from the start node.
	Walk(graph *DependencyGraph, startNode string) error
}

// NewWalker returns a Walker that will traverse the graph using a depth-first,
// post-order traversal.
func NewWalker(visitors ...DepVisitor) Walker {
	return &depthFirstWalker{
		visitors:   visitors,
		visitCache: make(map[string]int),
	}
}

type depthFirstWalker struct {

	// visitors are the actions to take on each dep visited in the traversal
	visitors []DepVisitor

	// visitCache maintains a mapping of the nodes visited
	visitCache map[string]int
}

// Walk starts a depth-first post-order traversal of the graph, starting at
// the Dependency with the given name.
func (w *depthFirstWalker) Walk(graph *DependencyGraph, startNode string) error {
	start := graph.Get(startNode)
	if start == nil {
		return fmt.Errorf("node %s not found", startNode)
	}
	return w.visit(graph, start)
}

// visit uses a nodeVisitor to attempt to visit the given Dependency,
// recursively visiting the Dependency's own dependencies by calling visit on
// them.
// TODO(nickt): detect cycles
func (w *depthFirstWalker) visit(graph *DependencyGraph, dep *Dependency) error {
	for _, v := range w.visitors {
		v.PreVisit(dep)
	}
	defer w.postVisit(dep)

	if dep == nil {
		return fmt.Errorf("dependency not found in graph")
	}

	// if we've visited this node, no need to visit again
	count, visited := w.visitCache[dep.Name]
	if visited {
		return nil
	}

	// however, if the count on this node is greater than zero, we're already visiting
	// this node higher up in the graph
	if count > 0 {
		return fmt.Errorf("detected cycle")
	}

	// and add this to the list of visited nodes
	w.visitCache[dep.Name]++

	// and visit all of the dependent nodes
	for _, d := range dep.Dependencies {
		dependency := graph.Get(d.Name)
		err := w.visit(graph, dependency)
		if err != nil {
			return fmt.Errorf("error visiting dependency '%s': %s", d.Name, err)
		}
	}

	// take action(s) on this current node
	for _, v := range w.visitors {
		if err := v.Visit(dep); err != nil {
			return err
		}
	}

	// "un-visit" this dep from the list of nodes
	w.visitCache[dep.Name]--

	return nil
}

func (w *depthFirstWalker) postVisit(dep *Dependency) {
	for _, v := range w.visitors {
		v.PostVisit(dep)
	}
}

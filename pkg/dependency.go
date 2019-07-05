package pkg

import (
	"fmt"

	"github.com/nicktrav/matryoshka/pkg/actions"
)

// state is the current state of the Dependency
type State int

const (
	unknown State = iota
	unsatisfied
	satisfied
)

func (s State) String() string {
	switch s {
	case unsatisfied:
		return "unsatisfied"
	case satisfied:
		return "satisfied"
	}
	return "unknown"
}

// A dependency represents a node in the dependency graph.
type Dependency struct {

	// Name is the name of the dependency.
	Name string

	// Dependencies is the list of dependencies that must be satisfied before
	// this dependency is satisfied.
	Dependencies []*Dependency

	// MetAction is the list of commands to run to determine whether the dependency is satisfied.
	// TODO(nickt): Change schema to be only one action
	MetActions []actions.Action

	// MeetAction is the list of commands to run to attempt to satisfy the dependency.
	MeetActions []actions.Action

	// State is the cached state of the Dependency.
	State
}

// NewDependency returns a pointer to a new Dependency.
func NewDependency(name string) *Dependency {
	return &Dependency{
		Name:         name,
		Dependencies: []*Dependency{},
		MetActions:   []actions.Action{},
		MeetActions:  []actions.Action{},
		State:        unknown,
	}
}

// Satisfy attempts to enforce the requirements of the Dependency, as defined
// by the MeetActions.
//
// A Dependency is considered "satisfied" if all of its dependencies are
// themselves satisfied, and the MetActions all return without error.
func (d *Dependency) Satisfy() (errors []error) {
	// if we're not in an unknown state, there's nothing to do.
	if d.State != unknown {
		return
	}

	// ensure that each dependency of the current dependency is satisfied
	// this check is optimistic in the sense that any dependency that fails
	// will not fail the check. All dependencies have a chance to be checked.
	for _, dep := range d.Dependencies {
		errs := dep.Satisfy()
		if len(errs) > 0 {
			d.State = unsatisfied
			errors = append(errors, errs...)
		}

		// if there were no errors, but the state of a dep was unsatisfied,
		// we're now also unsatisfied.
		if dep.State == unsatisfied {
			d.State = unsatisfied
		}
	}

	// if there were errors, there's nothing more that we can do, so we return
	// the errors that we've collected so far.
	if len(errors) > 0 {
		return
	}

	// if there were no errors, but the state was set to unsatisfied (due to
	// the checked deps having a cached state), we update our state and return.
	if d.State == unsatisfied {
		return
	}

	// check that all met actions are satisfied, aborting early if an action
	// fails.
	for _, a := range d.MetActions {
		if err := a.Run(); err != nil {
			d.State = unsatisfied
			break
		}
	}

	// if all of the above actions pass, we're satisfied and can return.
	if d.State != unsatisfied {
		d.State = satisfied
		return
	}

	// otherwise, we need to try and enforce the dependency by running the meet
	// actions. Any single action failure results in the dependency
	// transitioning to unsatisfied.
	for _, a := range d.MeetActions {
		if err := a.Run(); err != nil {
			d.State = unsatisfied
			errors = append(errors, fmt.Errorf("meet action: %s", err))
			return
		}
	}

	// if we completed the meet actions without error, we need to check again
	// that the desired state was indeed reached.
	for _, a := range d.MetActions {
		if err := a.Run(); err != nil {
			d.State = unsatisfied
			errors = append(errors, fmt.Errorf("met action (post-meet): %s", err))
			return
		}
	}

	d.State = satisfied
	return
}

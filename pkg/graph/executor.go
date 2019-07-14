package graph

import (
	"fmt"

	"github.com/nicktrav/matryoshka/pkg/actions"
)

// ExecutorOptions contains configuration information for an executor.
type ExecutorOptions struct {

	// Debug determines whether debug output will be printed.
	Debug bool

	// DryRun determines whether the "meet" action on each dep will be run.
	DryRun bool
}

// NewExecutor returns a new executor.
func NewExecutor(options ExecutorOptions) DepVisitor {
	return &executor{options}
}

// executor is a DepVisitor that will attempt to execute the "meet" actions on
// the Dependency it is visiting.
type executor struct {
	options ExecutorOptions
}

// Visit attempts to satisfy the current dependency, first checking that dep
// has not already been run, then checking the current deps deps to see if it
// is eligible to be run, then checking the current met actions. If the dep is
// still eligible to be run, the meet actions are run.
//
// This visitor assumes that a node is only visited once, hence the state of
// the dep should be unknown at the time of visiting.
func (e *executor) Visit(dep *Dependency) error {

	// we should have not visited this node before and the state should be unknown
	if dep.State != Unknown {
		return fmt.Errorf("executor: dep %s already visited", dep.Name)
	}

	// check actions of the deps below us
	for _, d := range dep.Dependencies {

		// if any are unsatisfied, we're also unsatisfied
		if d.State == Unsatisfied {
			dep.State = Unsatisfied
			return nil
		}
	}

	// else check the met actions of this node
	satisfied := true
	for _, a := range dep.MetActions {
		if e.options.Debug {
			enableDebug(a)
		}

		// if any of the met actions did not return successful, stop early
		if err := a.Run(); err != nil {
			satisfied = false
			break
		}
	}

	// if our met actions were all satisfied, this dep is satisfied
	if satisfied {
		dep.State = Satisfied
		return nil
	}

	// otherwise, we need to run the meet actions to attempt to enforce state
	for _, a := range dep.MeetActions {

		// if running in dry-run mode, don't run the action
		if e.options.DryRun {
			continue
		}

		if e.options.Debug {
			enableDebug(a)
		}

		// if any meet action could not be run, we're unsatisfied
		if err := a.Run(); err != nil {
			dep.State = Unsatisfied
			return nil
		}
	}

	// check the met actions again
	for _, a := range dep.MetActions {
		if e.options.Debug {
			enableDebug(a)
		}

		// if any of the met actions did not return successful, we're unsatisfied
		if err := a.Run(); err != nil {
			dep.State = Unsatisfied
			return nil
		}
	}

	// we made it through all the deps, this node is now satisfied
	dep.State = Satisfied
	return nil
}

// PreVisit does nothing.
func (e *executor) PreVisit(dep *Dependency) {
}

// PostVisit does nothing.
func (e *executor) PostVisit(dep *Dependency) {
}

// enableDebug enables debugging on the Action, if it is a Debugger.
func enableDebug(action actions.Action) {
	if debugger, ok := action.(actions.Debugger); ok {
		debugger.Debug()
	}
}

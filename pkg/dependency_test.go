package pkg

import (
	"errors"
	"fmt"
	"github.com/nicktrav/matryoshka/pkg/actions"
	"testing"
)

func TestDependency_Satisfy_StateKnown(t *testing.T) {
	dep := NewDependency("foo")

	action := countingAction{}
	dep.MetActions = []actions.Action{&action}
	dep.MeetActions = []actions.Action{&action}

	states := []State{satisfied, unsatisfied}
	for _, state := range states {
		dep.State = state
		errs := dep.Satisfy()
		if len(errs) > 0 {
			t.Errorf("wanted no errors; got %s", errs)
		}

		if action.count != 0 {
			t.Errorf("wanted action called 0 times; called %d times0", action.count)
		}
	}
}

func TestDependency_Satisfy_NoDepsOrActions(t *testing.T) {
	dep := NewDependency("foo")

	errs := dep.Satisfy()
	if len(errs) > 0 {
		t.Errorf("wanted no errors; got %s", errs)
	}

	if dep.State != satisfied {
		t.Errorf("wanted state satisfied; got %s", dep.State)
	}
}

func TestDependency_Satisfy_DepsCachedUnsatisfied(t *testing.T) {
	bar := NewDependency("bar")
	bar.State = unsatisfied

	foo := NewDependency("foo")
	foo.Dependencies = []*Dependency{bar}

	action := countingAction{}
	foo.MetActions = []actions.Action{&action}
	foo.MeetActions = []actions.Action{&action}

	errs := foo.Satisfy()
	if len(errs) > 0 {
		t.Errorf("wanted no errors; got %s", errs)
	}

	if foo.State != unsatisfied {
		t.Errorf("wanted state unsatisfied; got %s", foo.State)
	}

	if action.count > 0 {
		t.Errorf("wanted action called 0 times; called %d times", action.count)
	}
}

func TestDependency_Satisfy_DepsMeetActionFails(t *testing.T) {
	err := errors.New("failed")
	bar := NewDependency("bar")
	bar.MetActions = []actions.Action{&failingAction{err: err}}
	bar.MeetActions = []actions.Action{&failingAction{err: err}}

	foo := NewDependency("foo")
	foo.Dependencies = []*Dependency{bar}

	action := countingAction{}
	foo.MetActions = []actions.Action{&action}
	foo.MeetActions = []actions.Action{&action}

	errs := foo.Satisfy()
	if len(errs) != 1 {
		t.Errorf("wanted 1 error; got %d", len(errs))
	}

	want := fmt.Sprintf("meet action: %s", err)
	got := errs[0].Error()
	if want != got {
		t.Errorf("wanted %s; got %s", want, got)
	}

	if foo.State != unsatisfied {
		t.Errorf("wanted state unsatisfied; got %s", foo.State)
	}

	if action.count > 0 {
		t.Errorf("wanted action called 0 times; called %d times", action.count)
	}
}

func TestDependency_Satisfy_MetActionFailsThenSucceeds(t *testing.T) {
	dep := NewDependency("foo")

	metAction := failNTimesAction{n: 1, err: errors.New("fail")}
	dep.MetActions = []actions.Action{&metAction}

	meetAction := countingAction{}
	dep.MeetActions = []actions.Action{&meetAction}

	errs := dep.Satisfy()
	if len(errs) > 0 {
		t.Errorf("wanted no errors; got %s", errs)
	}

	if dep.State != satisfied {
		t.Errorf("wanted state satisfied; got %s", dep.State)
	}

	if metAction.count != 2 {
		t.Errorf("wanted metAction to be called twice; called %d times", metAction.count)
	}

	if meetAction.count != 1 {
		t.Errorf("wanted meetAction to be called once; called %d times", meetAction.count)
	}
}

func TestDependency_Satisfy_MetActionSucceeds(t *testing.T) {
	dep := NewDependency("foo")

	metAction := countingAction{}
	dep.MetActions = []actions.Action{&metAction}

	meetAction := countingAction{}
	dep.MeetActions = []actions.Action{&meetAction}

	errs := dep.Satisfy()
	if len(errs) > 0 {
		t.Errorf("wanted no errors; got %s", errs)
	}

	if dep.State != satisfied {
		t.Errorf("wanted state satisfied; got %s", dep.State)
	}

	if meetAction.count > 0 {
		t.Errorf("wanted action to not be called; called %d times", meetAction.count)
	}

	if metAction.count != 1 {
		t.Errorf("wanted met action to be called once; called %d times", metAction.count)
	}
}

func TestDependency_Satisfy_MetActionSucceeds_MeetActionFails(t *testing.T) {
	dep := NewDependency("foo")

	metAction := failingAction{err: errors.New("met action")}
	dep.MetActions = []actions.Action{&metAction}

	err := errors.New("fail")
	meetAction := failingAction{err: err}
	dep.MeetActions = []actions.Action{&meetAction}

	errs := dep.Satisfy()
	if len(errs) != 1 {
		t.Errorf("wanted one error; got %d", len(errs))
	}

	if dep.State != unsatisfied {
		t.Errorf("wanted state unsatisfied; got %s", dep.State)
	}

	want := fmt.Sprintf("meet action: %s", err)
	got := errs[0].Error()
	if want != got {
		t.Errorf("wanted %s; got %s", want, got)
	}

	if meetAction.count != 1 {
		t.Errorf("wanted meet action to be called once; called %d times", meetAction.count)
	}

	if metAction.count != 1 {
		t.Errorf("wanted met action to be called once; called %d times", metAction.count)
	}
}

// countingAction is an Action that counts the number of time is was called.
type countingAction struct {
	count int
}

func (a *countingAction) Run() error {
	a.count++
	return nil
}

// countingAction is an Action that fails with a given error, while recording
// the number of times it was called.
type failingAction struct {
	count int
	err   error
}

func (a *failingAction) Run() error {
	a.count++
	return a.err
}

// failNTimesAction is an Action that will fail a given number of times with an
// error, and then run without error.
type failNTimesAction struct {
	count int
	n     int
	err   error
}

func (a *failNTimesAction) Run() error {
	a.count++
	if a.count > a.n {
		return nil
	}
	return a.err
}

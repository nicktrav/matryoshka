package graph

import (
	"errors"
	"testing"

	"github.com/nicktrav/matryoshka/pkg/actions"
)

func TestExecutor_Visit_StateNotUnknown(t *testing.T) {
	deps := []*Dependency{
		{State: Unsatisfied},
		{State: Satisfied},
	}

	e := executor{}
	for _, dep := range deps {
		if err := e.Visit(dep); err == nil {
			t.Errorf("wanted error; got none")
		}
	}
}

func TestExecutor_Visit_DepsUnsatisfied(t *testing.T) {
	barDep := NewDependency("bar")
	barDep.State = Unsatisfied

	fooDep := NewDependency("foo")
	fooDep.Dependencies = []*Dependency{barDep}

	if fooDep.State != Unknown {
		t.Fatalf("wanted state unknown; got %+v", fooDep.State)
	}

	e := executor{}
	err := e.Visit(fooDep)
	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if fooDep.State != Unsatisfied {
		t.Fatalf("wanted state unsatisfied; got %+v", fooDep.State)
	}
}

func TestExecutor_Visit_DepsSatisfied_MetActionsSatisfied(t *testing.T) {
	barDep := NewDependency("bar")
	barDep.State = Satisfied

	fooDep := NewDependency("foo")
	fooDep.Dependencies = []*Dependency{barDep}

	metAction := &countingAction{}
	fooDep.MetActions = []actions.Action{metAction}

	meetAction := &countingAction{}
	fooDep.MeetActions = []actions.Action{meetAction}

	if fooDep.State != Unknown {
		t.Fatalf("wanted state unknown; got %+v", fooDep.State)
	}

	e := executor{}
	err := e.Visit(fooDep)
	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if fooDep.State != Satisfied {
		t.Fatalf("wanted state satisfied; got %+v", fooDep.State)
	}

	if metAction.count != 1 {
		t.Errorf("wanted met action called once; called %d times", metAction.count)
	}

	if meetAction.count != 0 {
		t.Errorf("wanted meet action called zero times; called %d times", meetAction.count)
	}
}

func TestExecutor_Visit_DepsSatisfied_MetActionsUnsatisfied_MeetActionsFail(t *testing.T) {
	barDep := NewDependency("bar")
	barDep.State = Satisfied

	fooDep := NewDependency("foo")
	fooDep.Dependencies = []*Dependency{barDep}

	metAction := newFailingAction()
	fooDep.MetActions = []actions.Action{metAction}

	meetAction := newFailingAction()
	fooDep.MeetActions = []actions.Action{meetAction}

	if fooDep.State != Unknown {
		t.Fatalf("wanted state unknown; got %+v", fooDep.State)
	}

	e := executor{}
	err := e.Visit(fooDep)
	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if fooDep.State != Unsatisfied {
		t.Fatalf("wanted state unsatisfied ; got %+v", fooDep.State)
	}

	if metAction.count != 1 {
		t.Errorf("wanted met action called once; called %d times", metAction.count)
	}

	if meetAction.count != 1 {
		t.Errorf("wanted meet action called zero times; called %d times", meetAction.count)
	}
}

func TestExecutor_Visit_DepsSatisfied_MetActionsUnsatisfied_MeetActionsSucceed_MetActionUnsatisfied(t *testing.T) {
	barDep := NewDependency("bar")
	barDep.State = Satisfied

	fooDep := NewDependency("foo")
	fooDep.Dependencies = []*Dependency{barDep}

	metAction := newFailingAction()
	fooDep.MetActions = []actions.Action{metAction}

	meetAction := &countingAction{}
	fooDep.MeetActions = []actions.Action{meetAction}

	if fooDep.State != Unknown {
		t.Fatalf("wanted state unknown; got %+v", fooDep.State)
	}

	e := executor{}
	err := e.Visit(fooDep)
	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if fooDep.State != Unsatisfied {
		t.Fatalf("wanted state unsatisfied ; got %+v", fooDep.State)
	}

	if metAction.count != 2 {
		t.Errorf("wanted met action called once; called %d times", metAction.count)
	}

	if meetAction.count != 1 {
		t.Errorf("wanted meet action called zero times; called %d times", meetAction.count)
	}
}

func TestExecutor_Visit_DepsSatisfied_MetActionsUnsatisfied_MeetActionsSucceed_MetActionSatisfied(t *testing.T) {
	barDep := NewDependency("bar")
	barDep.State = Satisfied

	fooDep := NewDependency("foo")
	fooDep.Dependencies = []*Dependency{barDep}

	metAction := &failNTimesAction{
		err: errors.New("oh noes"),
		n:   1, // fail once
	}
	fooDep.MetActions = []actions.Action{metAction}

	meetAction := &countingAction{}
	fooDep.MeetActions = []actions.Action{meetAction}

	if fooDep.State != Unknown {
		t.Fatalf("wanted state unknown; got %+v", fooDep.State)
	}

	e := executor{}
	err := e.Visit(fooDep)
	if err != nil {
		t.Fatalf("wanted no error; got %+v", err)
	}

	if fooDep.State != Satisfied {
		t.Fatalf("wanted state satisfied ; got %+v", fooDep.State)
	}

	if metAction.count != 2 {
		t.Errorf("wanted met action called once; called %d times", metAction.count)
	}

	if meetAction.count != 1 {
		t.Errorf("wanted meet action called zero times; called %d times", meetAction.count)
	}
}

func TestExecutor_EnableDebug_DebugAction(t *testing.T) {
	a := debugAction{}

	enableDebug(&a)

	if !a.debugCalled {
		t.Errorf("wanted debug called; debug not called")
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

func newFailingAction() *failingAction {
	return &failingAction{err: errors.New("oh noes")}
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

type debugAction struct {
	debugCalled bool
}

func (a *debugAction) Run() error {
	return nil
}

func (a *debugAction) Debug() {
	a.debugCalled = true
}

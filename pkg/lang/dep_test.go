package lang

import (
	"math/big"
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestDep_String(t *testing.T) {
	dep := Dep{Name: "foo"}

	want := "<dep.Dep \"foo\">"
	if want != dep.String() {
		t.Errorf("want %+v, got %+v", want, dep)
	}
}

func TestDep_Type(t *testing.T) {
	dep := Dep{}

	want := "dep.Dep"
	got := dep.Type()
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestDep_Truth(t *testing.T) {
	dep := Dep{}

	if !dep.Truth() {
		t.Error("Truth did not return true")
	}
}

func TestDep_Hash(t *testing.T) {
	dep := Dep{}

	_, err := dep.Hash()

	if err == nil {
		t.Fatal("wanted Hash to return an error")
	}

	if !strings.HasPrefix(err.Error(), "unhashable type") {
		t.Errorf("wanted error to contain 'unhashable type")
	}
}

func TestAsString(t *testing.T) {
	want := "foo"
	value := starlark.String(want)
	got, err := asString(value)

	if err != nil {
		t.Errorf("could not parse starlark string")
	}

	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestAsString_NotString(t *testing.T) {
	values := []starlark.Value{
		starlark.MakeInt(42),
		starlark.MakeBigInt(big.NewInt(42)),
		starlark.NewDict(42),
		starlark.NewList([]starlark.Value{}),
	}

	for _, value := range values {
		_, err := asString(value)
		if err == nil {
			t.Errorf("wanted error parsing value %+v to a string", value)
		}
	}
}

func TestAsRequirementsList(t *testing.T) {
	depItems := []*Dep{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}
	var values []starlark.Value
	for _, item := range depItems {
		values = append(values, item)
	}

	list, err := asRequirementsList(starlark.NewList(values))
	if err != nil {
		t.Errorf("error parsing requirements list: %s", err)
	}

	for i, item := range list {
		want := depItems[i]
		if want != item {
			t.Errorf("want %+v, got %+v", want, item)
		}
	}
}

func TestAsRequirementsList_NotDep(t *testing.T) {
	lists := []*starlark.List{
		starlark.NewList([]starlark.Value{starlark.String("foo")}),
		starlark.NewList([]starlark.Value{starlark.MakeInt(42)}),
		starlark.NewList([]starlark.Value{starlark.NewDict(42)}),
	}

	for _, list := range lists {
		_, err := asRequirementsList(list)
		if err == nil {
			t.Fatal("wanted error")
		}

		if !strings.Contains(err.Error(), "is not a dep") {
			t.Error("wanted error message to contain 'is not a dep'")
		}
	}
}

func TestAsRequirementsList_NotList(t *testing.T) {
	values := []starlark.Value{
		starlark.String("foo"),
		starlark.NewDict(42),
		starlark.MakeInt(42),
	}

	for _, value := range values {
		_, err := asRequirementsList(value)
		if err == nil {
			t.Fatalf("wanted error parsing %s as a list", value)
		}

		if !strings.Contains(err.Error(), "is not a list") {
			t.Errorf("wanted error message to contain 'is not a list'")
		}
	}
}

func TestAsCommandList(t *testing.T) {
	cmdItems := []ShellCmd{
		{Shell: "foo"},
		{Shell: "bar"},
		{Shell: "baz"},
	}
	var values []starlark.Value
	for _, item := range cmdItems {
		values = append(values, item)
	}

	list, err := asCommandList(starlark.NewList(values))
	if err != nil {
		t.Errorf("error parsing requirements list: %s", err)
	}

	for i, item := range list {
		want := cmdItems[i].Command
		got := item.Command
		if want != got {
			t.Errorf("want %+v, got %+v", want, got)
		}
	}
}

func TestAsCommandList_NotDep(t *testing.T) {
	lists := []*starlark.List{
		starlark.NewList([]starlark.Value{starlark.String("foo")}),
		starlark.NewList([]starlark.Value{starlark.MakeInt(42)}),
		starlark.NewList([]starlark.Value{starlark.NewDict(42)}),
	}

	for _, list := range lists {
		_, err := asCommandList(list)
		if err == nil {
			t.Fatal("wanted error")
		}

		if !strings.Contains(err.Error(), "is not a command") {
			t.Error("wanted error message to contain 'is not a command'")
		}
	}
}

func TestAsCommandList_NotList(t *testing.T) {
	values := []starlark.Value{
		starlark.String("foo"),
		starlark.NewDict(42),
		starlark.MakeInt(42),
	}

	for _, value := range values {
		_, err := asCommandList(value)
		if err == nil {
			t.Fatalf("wanted error parsing %s as a list", value)
		}

		if !strings.Contains(err.Error(), "is not a list") {
			t.Errorf("wanted error message to contain 'is not a list'")
		}
	}
}

func TestAsBool(t *testing.T) {
	type testCase struct {
		value          starlark.Value
		errorExpected  bool
		resultExpected bool
	}
	testCases := []testCase{
		{starlark.String("foo"), true, false},
		{starlark.Bool(true), false, true},
		{starlark.Bool(false), false, false},
	}

	for _, testCase := range testCases {
		res, err := asBool(testCase.value)

		if testCase.errorExpected && err == nil {
			t.Error("expected an error for input", testCase)
		}

		if !testCase.errorExpected && res != testCase.resultExpected {
			t.Error("expected", testCase.resultExpected, "got", res)
		}
	}
}
